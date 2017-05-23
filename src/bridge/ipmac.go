// Copyright (2013) Sandia Corporation.
// Under the terms of Contract DE-AC04-94AL85000 with Sandia Corporation,
// the U.S. Government retains certain rights in this software.

package bridge

import (
	log "minilog"
	"net"

)

func (b *Bridge) snooper() {
}

func (b *Bridge) updateIP(mac string, ip net.IP) {
	if ip == nil || ip.IsLinkLocalUnicast() {
		return
	}

	log.Debug("got mac/ip pair:", mac, ip)

	bridgeLock.Lock()
	defer bridgeLock.Unlock()

	for _, tap := range b.taps {
		if tap.Defunct || tap.MAC != mac {
			continue
		}

		if ip := ip.To4(); ip != nil {
			tap.IP4 = ip.String()
		} else {
			tap.IP6 = ip.String()
		}
	}
}
