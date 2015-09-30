package machine

import "github.com/stretchr/testify/mock"

type MockIdentifier struct {
	mock.Mock
}

func (self *MockIdentifier) Identify() ([]byte, error) {
	args := self.Called()

	obj := args.Get(0)

	if obj == nil {
		return nil, args.Error(1)
	}

	return obj.([]byte), args.Error(1)
}
