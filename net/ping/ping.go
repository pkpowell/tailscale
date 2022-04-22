// Copyright (c) 2022 Tailscale Inc & AUTHORS All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ping provides utilities for ICMP pings. It currently delegates to
// ping programs, in order to avoid capability and API restrictions.
package ping

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"

	"inet.af/netaddr"
	"tailscale.com/util/lineread"
	"tailscale.com/version/distro"
)

// ParseReply parses the first reply found in the reader r that matches the
// expected format of the current platform and returns the reported duration, or
// an error.
func ParseReply(r io.Reader) (time.Duration, netaddr.IP, error) {
	switch runtime.GOOS {
	case "windows":
		return parseReplyWindows(r)
	}
	return parseReplyUnix(r)
}

// Command constructs a single count ping command for the current platform targeting dstIP.
func Command(dstIP netaddr.IP) *exec.Cmd {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("ping", "-n", "1", "-w", "3000", dstIP.String())
	case "darwin":
		// Note: 2000 ms is actually 1 second + 2,000
		// milliseconds extra for 3 seconds total.
		// See https://github.com/tailscale/tailscale/pull/3753 for details.
		return exec.Command("ping", "-n", "-c", "1", "-W", "2000", dstIP.String())
	case "android":
		ping := "/system/bin/ping"
		if dstIP.Is6() {
			ping = "/system/bin/ping6"
		}
		return exec.Command(ping, "-n", "-c", "1", "-w", "3", dstIP.String())
	default:
		ping := "ping"
		if isSynology {
			ping = "/bin/ping"
		}
		cmd := exec.Command(ping, "-n", "-c", "1", "-W", "3", dstIP.String())
		if isSynology && os.Getuid() != 0 {
			// On DSM7 we run as non-root and need to pass
			// CAP_NET_RAW if our binary has it.
			setAmbientCapsRaw(cmd)
		}
		return cmd
	}
}

// setAmbientCapsRaw is non-nil on Linux for Synology, to run ping with
// CAP_NET_RAW from tailscaled's binary.
var setAmbientCapsRaw func(*exec.Cmd)

var isSynology = runtime.GOOS == "linux" && distro.Get() == distro.Synology

var stopErr = errors.New("stop")

// parseReplyWindows accepts an input stream of the form produced by the Windows
// ping utility, and parses the first reply line containing a duration,
// returning that duration. A reply that says <1ms is returned as
// time.Millisecond/2.
func parseReplyWindows(r io.Reader) (dur time.Duration, ip netaddr.IP, err error) {
	lrErr := lineread.Reader(r, func(line []byte) error {
		replyLine := bytes.TrimPrefix(line, []byte("Reply from "))
		if len(replyLine) == len(line) {
			return nil
		}

		endOfAddr := bytes.Index(replyLine, []byte(": "))
		if endOfAddr == -1 {
			return nil
		}
		addr := replyLine[:endOfAddr]

		ip, err = netaddr.ParseIP(string(addr))
		if err != nil {
			return stopErr
		}

		replyLine = replyLine[endOfAddr+2:]

		timeStart := bytes.Index(replyLine, []byte("time"))
		if timeStart == -1 {
			return nil
		}
		// the separator is either "=" or "<"
		timeSep := replyLine[timeStart+4]
		timeStart = timeStart + 5

		// the time=\d+ms may terminate the line, so matching on "ms".
		timeEnd := bytes.Index(replyLine, []byte("ms"))
		if timeEnd == -1 {
			return nil
		}

		// include the "ms", so ParseDuration knows the units
		timeB := replyLine[timeStart : timeEnd+2]

		dur, err = time.ParseDuration(string(timeB))
		if err != nil {
			return stopErr
		}

		if err == nil && timeSep == byte('<') {
			dur = dur / 2
		}

		return stopErr
	})
	if lrErr != stopErr {
		err = lrErr
	}
	return
}

// parseReplyUnix accepts an input stream of the form produced by the Unix
// ping utility, and parses the first reply line containing a duration,
// returning that duration.
func parseReplyUnix(r io.Reader) (dur time.Duration, ip netaddr.IP, err error) {
	lrErr := lineread.Reader(r, func(line []byte) error {
		replyPrefix := []byte(" bytes from ")

		prefixIdx := bytes.Index(line, replyPrefix)
		if prefixIdx == -1 {
			return nil
		}
		addrEtc := line[prefixIdx+len(replyPrefix):]

		addrEnd := bytes.Index(addrEtc, []byte(": "))
		if addrEnd == -1 {
			return nil
		}

		ip, err = netaddr.ParseIP(string(addrEtc[:addrEnd]))
		if err != nil {
			return stopErr
		}

		replyLine := addrEtc[addrEnd+2:]

		timeStart := bytes.Index(replyLine, []byte("time="))
		if timeStart == -1 {
			return nil
		}
		timeStart = timeStart + 5

		timeEnd := bytes.Index(replyLine, []byte(" ms"))
		if timeEnd == -1 {
			return nil
		}

		dur, err = time.ParseDuration(string(replyLine[timeStart:timeEnd]) + "ms")
		return stopErr
	})
	if lrErr != stopErr {
		err = lrErr
	}
	return
}
