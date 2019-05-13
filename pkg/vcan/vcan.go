package vcan

import (
	"github.com/vishvananda/netlink"
)

type Vcan struct {
	netlink.LinkAttrs
}

func (vcan *Vcan) Attrs() *netlink.LinkAttrs {
	return &vcan.LinkAttrs
}
func (vcan *Vcan) Type() string {
	return "vcan"
}
