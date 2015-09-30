package machine

import (
	"errors"
	"testing"
)

func TestChainingIdentifierCallsIntoCurrent(t *testing.T) {
	current := &MockIdentifier{}
	next := &MockIdentifier{}

	current.On("Identify").Return([]byte{42, 42, 42}, nil)

	ci := ChainingIdentifier{current, next}

	ci.Identify()

	current.AssertExpectations(t)
}

func TestChainingIdentifierCallsIntoNextOnError(t *testing.T) {
	current := &MockIdentifier{}
	next := &MockIdentifier{}

	current.On("Identify").Return(nil, errors.New("test"))
	next.On("Identify").Return([]byte{42, 42, 42}, nil)

	ci := ChainingIdentifier{current, next}

	ci.Identify()

	current.AssertExpectations(t)
	next.AssertExpectations(t)
}
