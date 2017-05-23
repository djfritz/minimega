// Copyright (2015) Sandia Corporation.
// Under the terms of Contract DE-AC04-94AL85000 with Sandia Corporation,
// the U.S. Government retains certain rights in this software.

package minipager

import (
	log "minilog"
)

type Pager interface {
	Page(output string)
}

// Copy of winsize struct defined by iotctl.h
type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

var DefaultPager Pager = &defaultPager{}

type defaultPager struct{}

func (_ defaultPager) Page(output string) {
	log.Fatalln("not implemented")
}
