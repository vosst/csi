package pkg

// System models a packaging system
type System interface {
	Resolver             // System provides means to resolve bundles given a search pattern
	Arch() (Arch, error) // Arch returns the machine architecture that the current system was built for
}
