package main

import (
	"github.com/vosst/whoopsie/nm"
	"launchpad.net/go-dbus/v1"
)

type NMReachabilityMonitor struct {
	ch             chan Reachability
	networkManager *nm.Manager
}

func NewNMReachabilityMonitor(conn *dbus.Connection) (*NMReachabilityMonitor, error) {
	if manager, err := nm.NewManager(conn); err != nil {
		return nil, err
	} else {
		return &NMReachabilityMonitor{make(chan Reachability), manager}, nil
	}
}
