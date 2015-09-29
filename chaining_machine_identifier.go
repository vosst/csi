package main

type ChainingMachineIdentifier struct {
	Identifier MachineIdentifier
	Next       MachineIdentifier
}

func (self ChainingMachineIdentifier) Identify() ([]byte, error) {
	b, err := self.Identifier.Identify()

	if err != nil {
		return self.Next.Identify()
	}

	return b, err
}
