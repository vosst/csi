package main

import "github.com/stretchr/testify/mock"

type MockReachabilityMonitor struct {
	mock.Mock
}

func (self *MockReachabilityMonitor) CheckHostReachability(host string) Reachability {
	args := self.Called(host)
	return Reachability(args.Int(0))
}
