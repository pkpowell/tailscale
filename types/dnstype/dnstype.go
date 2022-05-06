// Copyright (c) 2021 Tailscale Inc & AUTHORS All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package dnstype defines types for working with DNS.
package dnstype

//go:generate go run tailscale.com/cmd/cloner --type=Resolver --clonefunc=true

import "inet.af/netaddr"

// Resolver is the configuration for one DNS resolver.
type Resolver struct {
	// Addr is the address of the DNS resolver, one of:
	//  - A plain IP address for a "classic" UDP+TCP DNS resolver.
	//    This is the common format as sent by the control plane.
	//  - An IP:port, for tests.
	//  - [TODO] "tls://resolver.com" for DNS over TCP+TLS
	//  - [TODO] "https://resolver.com/query-tmpl" for DNS over HTTPS
	Addr string `json:",omitempty"`

	// BootstrapResolution is an optional suggested resolution for the
	// DoT/DoH resolver, if the resolver URL does not reference an IP
	// address directly.
	// BootstrapResolution may be empty, in which case clients should
	// look up the DoT/DoH server using their local "classic" DNS
	// resolver.
	BootstrapResolution []netaddr.IP `json:",omitempty"`
}

// IPPort returns r.Addr as an IP address and port if either
// r.Addr is an IP address (the common case) or if r.Addr
// is an IP:port (as done in tests).
func (r *Resolver) IPPort() (ipp netaddr.IPPort, ok bool) {
	if r.Addr == "" || r.Addr[0] == 'h' || r.Addr[0] == 't' {
		// Fast path to avoid ParseIP error allocation for obviously not IP
		// cases.
		return
	}
	if ip, err := netaddr.ParseIP(r.Addr); err == nil {
		return netaddr.IPPortFrom(ip, 53), true
	}
	if ipp, err := netaddr.ParseIPPort(r.Addr); err == nil {
		return ipp, true
	}
	return
}
