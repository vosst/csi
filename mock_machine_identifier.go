package main

import "github.com/stretchr/testify/mock"

type MockMachineIdentifier struct {
	mock.Mock
}

func (self *MockMachineIdentifier) Identify() ([]byte, error) {
	args := self.Called()
	return args.Get(0).([]byte), args.Error(1)
}
