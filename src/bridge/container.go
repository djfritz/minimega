// Copyright (2016) Sandia Corporation.
// Under the terms of Contract DE-AC04-94AL85000 with Sandia Corporation,
// the U.S. Government retains certain rights in this software.

package bridge

import (
	log "minilog"
)

// CreateContainerTap creates and adds a veth tap, t, to a bridge. If a name is
// not provided, one will be automatically generated. pid is the network
// namespace pid for the tap. mac is the MAC address to assign to the interface.
// vlan is the VLAN for the traffic. index is the veth interface number for the
// container.
func (b *Bridge) CreateContainerTap(t string, mac string, vlan, index int, pid int) (tap string, err error) {
	bridgeLock.Lock()
	defer bridgeLock.Unlock()

	tap = t
	if tap == "" {
		tap = <-b.nameChan
	}

	if err := createVeth(tap, pid, mac, index); err != nil {
		return "", err
	}

	// Clean up the tap we just created, if it didn't already exist.
	defer func() {
		if err != nil {
			if err := destroyTap(tap); err != nil {
				// Welp, we're boned
				log.Error("zombie tap -- %v %v", tap, err)
			}
			tap = ""
		}
	}()

	if err = upInterface(tap, false); err != nil {
		return "", err
	}

	// Add the interface
	return tap, b.addTap(tap, mac, vlan, false)
}
