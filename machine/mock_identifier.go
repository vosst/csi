package machine

import "github.com/stretchr/testify/mock"

type MockIdentifier struct {
	mock.Mock
}

func (self *MockIdentifier) Identify() ([]byte, error) {
	args := self.Called()
	return args.Get(0).([]byte), args.Error(1)
}
