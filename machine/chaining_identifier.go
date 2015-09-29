package machine

type ChainingIdentifier struct {
	Identifier Identifier
	Next       Identifier
}

func (self ChainingIdentifier) Identify() ([]byte, error) {
	b, err := self.Identifier.Identify()

	if err != nil {
		return self.Next.Identify()
	}

	return b, err
}
