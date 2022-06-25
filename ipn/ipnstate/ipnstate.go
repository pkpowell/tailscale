// Copyright (c) 2020 Tailscale Inc & AUTHORS All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ipnstate captures the entire state of the Tailscale network.
//
// It's a leaf package so ipn, wgengine, and magicsock can all depend on it.
package ipnstate

import (
	"embed"
	"fmt"
	"html"
	"log"
	"net"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"inet.af/netaddr"
	"tailscale.com/tailcfg"
	"tailscale.com/types/key"
	"tailscale.com/types/views"
	"tailscale.com/util/dnsname"
)

// Status represents the entire state of the IPN network.
type Status struct {
	// Version is the daemon's long version (see version.Long).
	Version string

	// BackendState is an ipn.State string value:
	//  "NoState", "NeedsLogin", "NeedsMachineAuth", "Stopped",
	//  "Starting", "Running".
	BackendState string

	AuthURL      string       // current URL provided by control to authorize client
	TailscaleIPs []netaddr.IP // Tailscale IP(s) assigned to this node
	Self         *PeerStatus

	// ExitNodeStatus describes the current exit node.
	// If nil, an exit node is not in use.
	ExitNodeStatus *ExitNodeStatus `json:"ExitNodeStatus,omitempty"`

	// Health contains health check problems.
	// Empty means everything is good. (or at least that no known
	// problems are detected)
	Health []string

	// This field is the legacy name of CurrentTailnet.MagicDNSSuffix.
	//
	// Deprecated: use CurrentTailnet.MagicDNSSuffix instead.
	MagicDNSSuffix string

	// CurrentTailnet is information about the tailnet that the node
	// is currently connected to. When not connected, this field is nil.
	CurrentTailnet *TailnetStatus

	// CertDomains are the set of DNS names for which the control
	// plane server will assist with provisioning TLS
	// certificates. See SetDNSRequest for dns-01 ACME challenges
	// for e.g. LetsEncrypt. These names are FQDNs without
	// trailing periods, and without any "_acme-challenge." prefix.
	CertDomains []string

	Peer map[key.NodePublic]*PeerStatus
	User map[tailcfg.UserID]tailcfg.UserProfile
}

// TailnetStatus is information about a Tailscale network ("tailnet").
type TailnetStatus struct {
	// Name is the name of the network that's currently in use.
	Name string

	// MagicDNSSuffix is the network's MagicDNS suffix for nodes
	// in the network such as "userfoo.tailscale.net".
	// There are no surrounding dots.
	// MagicDNSSuffix should be populated regardless of whether a domain
	// has MagicDNS enabled.
	MagicDNSSuffix string

	// MagicDNSEnabled is whether or not the network has MagicDNS enabled.
	// Note that the current device may still not support MagicDNS if
	// `--accept-dns=false` was used.
	MagicDNSEnabled bool
}

// ExitNodeStatus describes the current exit node.
type ExitNodeStatus struct {
	// ID is the exit node's ID.
	ID tailcfg.StableNodeID

	// Online is whether the exit node is alive.
	Online bool

	// TailscaleIPs are the exit node's IP addresses assigned to the node.
	TailscaleIPs []netaddr.IPPrefix
}

func (s *Status) Peers() []key.NodePublic {
	kk := make([]key.NodePublic, 0, len(s.Peer))
	for k := range s.Peer {
		kk = append(kk, k)
	}
	sort.Slice(kk, func(i, j int) bool { return kk[i].Less(kk[j]) })
	return kk
}

type PeerStatusLite struct {
	// TxBytes/RxBytes is the total number of bytes transmitted to/received from this peer.
	TxBytes, RxBytes int64
	// LastHandshake is the last time a handshake succeeded with this peer.
	// (Or we got key confirmation via the first data message,
	// which is approximately the same thing.)
	LastHandshake time.Time
	// NodeKey is this peer's public node key.
	NodeKey key.NodePublic
}

type PeerStatus struct {
	ID           tailcfg.StableNodeID
	PublicKey    key.NodePublic
	HostName     string // HostInfo's Hostname (not a DNS name or necessarily unique)
	DNSName      string
	OS           string // HostInfo.OS
	UserID       tailcfg.UserID
	TailscaleIPs []netaddr.IP // Tailscale IP(s) assigned to this node
	PeerAPIPort  uint16

	// Tags are the list of ACL tags applied to this node.
	// See tailscale.com/tailcfg#Node.Tags for more information.
	Tags *views.Slice[string] `json:",omitempty"`

	// PrimaryRoutes are the routes this node is currently the primary
	// subnet router for, as determined by the control plane. It does
	// not include the IPs in TailscaleIPs.
	PrimaryRoutes *views.IPPrefixSlice `json:",omitempty"`

	// Endpoints:
	Addrs   []string
	CurAddr string // one of Addrs, or unique if roaming
	Relay   string // DERP region

	RxBytes        int64
	TxBytes        int64
	Created        time.Time // time registered with tailcontrol
	LastWrite      time.Time // time last packet sent
	LastSeen       time.Time // last seen to tailcontrol; only present if offline
	LastHandshake  time.Time // with local wireguard
	Online         bool      // whether node is connected to the control plane
	KeepAlive      bool
	ExitNode       bool // true if this is the currently selected exit node.
	ExitNodeOption bool // true if this node can be an exit node (offered && approved)

	// Active is whether the node was recently active. The
	// definition is somewhat undefined but has historically and
	// currently means that there was some packet sent to this
	// peer in the past two minutes. That definition is subject to
	// change.
	Active bool

	PeerAPIURL   []string
	Capabilities []string `json:",omitempty"`

	// SSH_HostKeys are the node's SSH host keys, if known.
	SSH_HostKeys []string `json:"sshHostKeys,omitempty"`

	// ShareeNode indicates this node exists in the netmap because
	// it's owned by a shared-to user and that node might connect
	// to us. These nodes should be hidden by "tailscale status"
	// etc by default.
	ShareeNode bool `json:",omitempty"`

	// InNetworkMap means that this peer was seen in our latest network map.
	// In theory, all of InNetworkMap and InMagicSock and InEngine should all be true.
	InNetworkMap bool

	// InMagicSock means that this peer is being tracked by magicsock.
	// In theory, all of InNetworkMap and InMagicSock and InEngine should all be true.
	InMagicSock bool

	// InEngine means that this peer is tracked by the wireguard engine.
	// In theory, all of InNetworkMap and InMagicSock and InEngine should all be true.
	InEngine bool
}

type StatusBuilder struct {
	mu     sync.Mutex
	locked bool
	st     Status
}

// MutateStatus calls f with the status to mutate.
//
// It may not assume other fields of status are already populated, and
// may not retain or write to the Status after f returns.
//
// MutateStatus acquires a lock so f must not call back into sb.
func (sb *StatusBuilder) MutateStatus(f func(*Status)) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	f(&sb.st)
}

func (sb *StatusBuilder) Status() *Status {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	sb.locked = true
	return &sb.st
}

// MutateSelfStatus calls f with the PeerStatus of our own node to mutate.
//
// It may not assume other fields of status are already populated, and
// may not retain or write to the Status after f returns.
//
// MutateStatus acquires a lock so f must not call back into sb.
func (sb *StatusBuilder) MutateSelfStatus(f func(*PeerStatus)) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	if sb.st.Self == nil {
		sb.st.Self = new(PeerStatus)
	}
	f(sb.st.Self)
}

// AddUser adds a user profile to the status.
func (sb *StatusBuilder) AddUser(id tailcfg.UserID, up tailcfg.UserProfile) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	if sb.locked {
		log.Printf("[unexpected] ipnstate: AddUser after Locked")
		return
	}

	if sb.st.User == nil {
		sb.st.User = make(map[tailcfg.UserID]tailcfg.UserProfile)
	}

	sb.st.User[id] = up
}

// AddIP adds a Tailscale IP address to the status.
func (sb *StatusBuilder) AddTailscaleIP(ip netaddr.IP) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	if sb.locked {
		log.Printf("[unexpected] ipnstate: AddIP after Locked")
		return
	}

	sb.st.TailscaleIPs = append(sb.st.TailscaleIPs, ip)
}

// AddPeer adds a peer node to the status.
//
// Its PeerStatus is mixed with any previous status already added.
func (sb *StatusBuilder) AddPeer(peer key.NodePublic, st *PeerStatus) (e *PeerStatus) {

	if st == nil {
		panic("nil PeerStatus")
	}

	sb.mu.Lock()
	defer sb.mu.Unlock()
	if sb.locked {
		log.Printf("[unexpected] ipnstate: AddPeer after Locked")
		return
	}

	if sb.st.Peer == nil {
		sb.st.Peer = make(map[key.NodePublic]*PeerStatus)
	}
	e, ok := sb.st.Peer[peer]
	if !ok {
		sb.st.Peer[peer] = st
		st.PublicKey = peer
		return
	}

	// log.Printf("adding peer %s", e.HostName)

	if v := st.ID; v != "" {
		e.ID = v
	}
	if v := st.HostName; v != "" {
		e.HostName = v
	}
	if v := st.DNSName; v != "" {
		e.DNSName = v
	}
	if v := st.Relay; v != "" {
		e.Relay = v
	}
	if v := st.UserID; v != 0 {
		e.UserID = v
	}
	if v := st.TailscaleIPs; v != nil {
		e.TailscaleIPs = v
	}
	if v := st.PrimaryRoutes; v != nil && !v.IsNil() {
		e.PrimaryRoutes = v
	}
	if v := st.Tags; v != nil && !v.IsNil() {
		e.Tags = v
	}
	if v := st.OS; v != "" {
		e.OS = st.OS
	}
	if v := st.SSH_HostKeys; v != nil {
		e.SSH_HostKeys = v
	}
	if v := st.Addrs; v != nil {
		e.Addrs = v
	}
	if v := st.CurAddr; v != "" {
		e.CurAddr = v
	}
	if v := st.RxBytes; v != 0 {
		e.RxBytes = v
	}
	if v := st.TxBytes; v != 0 {
		e.TxBytes = v
	}
	if v := st.LastHandshake; !v.IsZero() {
		e.LastHandshake = v
	}
	if v := st.Created; !v.IsZero() {
		e.Created = v
	}
	if v := st.LastSeen; !v.IsZero() {
		e.LastSeen = v
	}
	if v := st.LastWrite; !v.IsZero() {
		e.LastWrite = v
	}
	if st.Online {
		e.Online = true
	}
	if st.InNetworkMap {
		e.InNetworkMap = true
	}
	if st.InMagicSock {
		e.InMagicSock = true
	}
	if st.InEngine {
		e.InEngine = true
	}
	if st.KeepAlive {
		e.KeepAlive = true
	}
	if st.ExitNode {
		e.ExitNode = true
	}
	if st.ExitNodeOption {
		e.ExitNodeOption = true
	}
	if st.ShareeNode {
		e.ShareeNode = true
	}
	if st.Active {
		e.Active = true
	}

	return e
}

type StatusUpdater interface {
	UpdateStatus(*StatusBuilder)
}

type statusData struct {
	Peers          []PeerData
	Profile        tailcfg.UserProfile
	DeviceName     string
	BackendState   string
	Version        string
	CurrentTailnet *TailnetStatus
	Health         []string
	IPs            []string
	IPv4           string
	IPv6           string
	Now            time.Time
}

// PeerData struct to export data to html template
type PeerData struct {
	IPs         []netaddr.IP         `json:"IPs"`
	IPv4        string               `json:"IPv4"`
	IPv6        string               `json:"IPv6"`
	IPv4Num     string               `json:"IPv4Num"`
	IPv6Num     string               `json:"IPv6Num"`
	Peer        string               `json:"Peer"`
	ActAgo      string               `json:"ActAgo"`
	LastSeen    string               `json:"LastSeen"`
	Unseen      bool                 `json:"Unseen"`
	Created     time.Time            `json:"Created"`
	CreatedDate string               `json:"CreatedDate"`
	CreatedTime string               `json:"CreatedTime"`
	Active      bool                 `json:"Active"`
	Online      bool                 `json:"Online"`
	OS          string               `json:"OS"`
	HostName    string               `json:"HostName"`
	HostInfo    tailcfg.Hostinfo     `json:"HostInfo"`
	ID          tailcfg.StableNodeID `json:"ID"`
	UserID      tailcfg.UserID       `json:"UserID"`
	NodeKey     string               `json:"NodeKey"`
	Owner       string               `json:"Owner"`
	DNSName     string               `json:"DNSName"`
	TailAddr    []net.IPAddr         `json:"TailAddr"`
	Connection  string               `json:"Connection"`
	TX          string               `json:"TX"`
	RX          string               `json:"RX"`
	TXb         int64                `json:"TXb"`
	RXb         int64                `json:"RXb"`
}

//go:embed assets/*.html assets/*.css
var assets embed.FS

func (st *Status) WriteHTMLtmpl(w http.ResponseWriter) {
	var data statusData
	data.Profile = st.User[st.Self.UserID]
	data.Now = time.Now()
	data.Version = st.Version
	data.Health = st.Health
	data.BackendState = st.BackendState
	data.CurrentTailnet = st.CurrentTailnet
	data.DeviceName = strings.Split(st.Self.DNSName, ".")[0]

	var peers []*PeerStatus
	for _, peer := range st.Peers() {
		ps := st.Peer[peer]
		if ps.ShareeNode {
			continue
		}
		peers = append(peers, ps)
	}
	SortPeers(peers)
	data.Peers = make([]PeerData, len(peers))

	data.IPs = make([]string, 0, len(st.TailscaleIPs))
	for _, ip := range st.TailscaleIPs {
		if ip.Is4() {
			data.IPv4 = ip.String()
		} else {
			data.IPv6 = ip.String()
		}
	}

	for i, ps := range peers {
		data.Peers[i].ID = ps.ID
		data.Peers[i].Peer = ps.PublicKey.ShortString()
		data.Peers[i].HostName = ps.HostName
		data.Peers[i].OS = ps.OS
		data.Peers[i].Online = ps.Online
		data.Peers[i].Active = ps.Active
		data.Peers[i].HostName = dnsname.SanitizeHostname(ps.HostName)
		data.Peers[i].IPs = ps.TailscaleIPs
		for _, ip := range ps.TailscaleIPs {
			if ip.Is4() {
				data.Peers[i].IPv4 = ip.String()
			} else {
				data.Peers[i].IPv6 = ip.String()
			}
		}
		if !ps.LastWrite.IsZero() {
			ago := data.Now.Sub(ps.LastWrite)
			data.Peers[i].ActAgo = ago.Round(time.Second).String() + " ago"
		}

		if up, ok := st.User[ps.UserID]; ok {
			data.Peers[i].Owner = up.LoginName
			if i := strings.Index(data.Peers[i].Owner, "@"); i != -1 {
				data.Peers[i].Owner = data.Peers[i].Owner[:i]
			}
		}

		data.Peers[i].DNSName = dnsname.TrimSuffix(ps.DNSName, st.MagicDNSSuffix)

		data.Peers[i].RX = FormatBytes(ps.RxBytes, Base2)
		data.Peers[i].TX = FormatBytes(ps.TxBytes, Base2)

		if ps.Active {
			if ps.Relay != "" && ps.CurAddr == "" {
				data.Peers[i].Connection = html.EscapeString(ps.Relay)
				// } else if ps.CurAddr != "" {
				// 	data.Peers[i].Connection = html.EscapeString(ps.CurAddr)
			}
		}
	}

	// buf := new(bytes.Buffer)
	// if err := tmpl.Execute(buf, data); err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	fmt.Printf("an error happened %v", err)
	// }
	// w.Write(buf.Bytes())
}

type Prefix struct {
	full  string
	short byte
}

type Unit struct {
	factor int64
	suffix string
}

var (
	Base2  = Unit{factor: 1024, suffix: "iB"}
	Base10 = Unit{factor: 1000, suffix: "B"}

	PRX = map[int]Prefix{
		0: {
			short: 'K',
			full:  "kilo",
		},
		1: {
			short: 'M',
			full:  "mega",
		},
		2: {
			short: 'G',
			full:  "giga",
		},
		3: {
			short: 'T',
			full:  "tera",
		},
		4: {
			short: 'P',
			full:  "peta",
		},
		5: {
			short: 'E',
			full:  "exa",
		},
		// these are just for show. int64 only reaches 8EB
		6: {
			short: 'Y',
			full:  "yotta",
		},
		7: {
			short: 'Z',
			full:  "zetta",
		},
	}
)

// FormatBytes converts bytes to KB, MiB etc without Math lib
func FormatBytes(b int64, u Unit) string {
	// bytes only
	if b == 0 {
		return "–"
	}
	if b < u.factor {
		return fmt.Sprintf("%d B", b)
	}

	div := u.factor
	exp := 0

	for n := b / u.factor; n >= u.factor; n /= u.factor {
		// grow the divisor
		div *= u.factor
		exp++
	}
	return fmt.Sprintf("%.2f %c%s", float64(b)/float64(div), PRX[exp].short, u.suffix)
}

// PingResult contains response information for the "tailscale ping" subcommand,
// saying how Tailscale can reach a Tailscale IP or subnet-routed IP.
// See tailcfg.PingResponse for a related response that is sent back to control
// for remote diagnostic pings.
type PingResult struct {
	IP       string // ping destination
	NodeIP   string // Tailscale IP of node handling IP (different for subnet routers)
	NodeName string // DNS name base or (possibly not unique) hostname

	Err            string
	LatencySeconds float64

	// Endpoint is the ip:port if direct UDP was used.
	// It is not currently set for TSMP pings.
	Endpoint string

	// DERPRegionID is non-zero DERP region ID if DERP was used.
	// It is not currently set for TSMP pings.
	DERPRegionID int

	// DERPRegionCode is the three-letter region code
	// corresponding to DERPRegionID.
	// It is not currently set for TSMP pings.
	DERPRegionCode string

	// PeerAPIPort is set by TSMP ping responses for peers that
	// are running a peerapi server. This is the port they're
	// running the server on.
	PeerAPIPort uint16 `json:",omitempty"`

	// PeerAPIURL is the URL that was hit for pings of type "peerapi" (tailcfg.PingPeerAPI).
	// It's of the form "http://ip:port" (or [ip]:port for IPv6).
	PeerAPIURL string `json:",omitempty"`

	// IsLocalIP is whether the ping request error is due to it being
	// a ping to the local node.
	IsLocalIP bool `json:",omitempty"`

	// TODO(bradfitz): details like whether port mapping was used on either side? (Once supported)
}

func (pr *PingResult) ToPingResponse(pingType tailcfg.PingType) *tailcfg.PingResponse {
	return &tailcfg.PingResponse{
		Type:           pingType,
		IP:             pr.IP,
		NodeIP:         pr.NodeIP,
		NodeName:       pr.NodeName,
		Err:            pr.Err,
		LatencySeconds: pr.LatencySeconds,
		Endpoint:       pr.Endpoint,
		DERPRegionID:   pr.DERPRegionID,
		DERPRegionCode: pr.DERPRegionCode,
		PeerAPIPort:    pr.PeerAPIPort,
		IsLocalIP:      pr.IsLocalIP,
	}
}

func SortPeers(peers []*PeerStatus) {
	sort.Slice(peers, func(i, j int) bool { return sortKey(peers[i]) < sortKey(peers[j]) })
}

func sortKey(ps *PeerStatus) string {
	if ps.DNSName != "" {
		return ps.DNSName
	}
	if ps.HostName != "" {
		return ps.HostName
	}
	// TODO(bradfitz): add PeerStatus.Less and avoid these allocs in a Less func.
	if len(ps.TailscaleIPs) > 0 {
		return ps.TailscaleIPs[0].String()
	}
	raw := ps.PublicKey.Raw32()
	return string(raw[:])
}
