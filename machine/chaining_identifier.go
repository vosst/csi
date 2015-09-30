package machine

// ChainingIdentifier satisfies Identify requests relying
// on a current and a next Identifier. If the call to Current
// fails, the request is handed over to Next.
type ChainingIdentifier struct {
	Current Identifier
	Next    Identifier
}

// Identify establishes an identity by first invoking Current.
// If Current.Identify fails, it hands over to Next.
func (self ChainingIdentifier) Identify() ([]byte, error) {
	b, err := self.Current.Identify()

	if err != nil {
		return self.Next.Identify()
	}

	return b, err
}
