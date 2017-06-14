// Copyright (2016) Sandia Corporation.
// Under the terms of Contract DE-AC04-94AL85000 with Sandia Corporation,
// the U.S. Government retains certain rights in this software.

package bridge

import (
	"fmt"
	log "minilog"
	"net"

	"github.com/vishvananda/netlink"
)

// upInterface activates an interface using the `ip` command. promisc controls
// whether the interface is brought up in promiscuous mode or not.
func upInterface(name string, promisc bool) error {
	log.Info("up interface: %v", name)

	la := netlink.NewLinkAttrs()
	la.Name = name
	dev := &netlink.Device{
		LinkAttrs: la,
	}

	err := netlink.LinkSetUp(dev)
	if err != nil {
		return err
	}

	if promisc {
		return netlink.SetPromiscOn(dev)
	}

	return nil
}

// downInterface deactivates an interface using the `ip` command.
func downInterface(name string) error {
	log.Info("down interface: %v", name)

	la := netlink.NewLinkAttrs()
	la.Name = name
	dev := &netlink.Device{
		LinkAttrs: la,
	}

	return netlink.LinkSetDown(dev)
}

// createTap creates a tuntap of the specified name using the `ip` command.
func createTap(name string) error {
	log.Info("creating tuntap: %v", name)

	la := netlink.NewLinkAttrs()
	la.Name = name
	dev := &netlink.Tuntap{
		LinkAttrs: la,
	}

	return netlink.LinkAdd(dev)
}

// createVeth creates a veth of the specified name using the `ip` command.
func createVeth(name string, pid int, mac string, index int) error {
	log.Debug("creating veth: %v", name)

	ifname := fmt.Sprintf("veth%v", index)

	hwaddr, err := net.ParseMAC(mac)
	if err != nil {
		return err
	}

	la := netlink.NewLinkAttrs()
	la.Name = ifname
	la.Namespace = netlink.NsPid(pid)
	la.HardwareAddr = hwaddr
	veth := &netlink.Veth{
		LinkAttrs: la,
		PeerName:  name,
	}
	return netlink.LinkAdd(veth)
}

// DestroyTap destroys an `unmanaged` tap using the `ip` command. This can be
// used when cleaning up from a crash or when a tap is not connected to a
// bridges. See `Bridge.DestroyTap` for managed taps.
func DestroyTap(name string) error {
	bridgeLock.Lock()
	defer bridgeLock.Unlock()

	return destroyTap(name)
}

// destroyTap destroys a tuntap device.
func destroyTap(name string) error {
	log.Info("destroying tuntap: %v", name)

	la := netlink.NewLinkAttrs()
	la.Name = name
	dev := &netlink.Device{
		LinkAttrs: la,
	}

	return netlink.LinkDel(dev)
}
