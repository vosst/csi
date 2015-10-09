package pkg

// Bundle describes a bundle, for example a .snap, installed on the system.
type Bundle interface {
	// Name returns the name of the bundle
	Name() string
	// Version returns the version of the bundle
	Version() string
	// Arch returns the architecture that the bundle has been built for
	Arch() string
}
