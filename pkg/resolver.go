package pkg

// Resolver abstracts resolving of Bundles given a pattern or filename.
type Resolver interface {
	// Resolve returns all Bundles matching pattern.
	//
	// Returns an error if querying the underlying index fails.
	Resolve(pattern string) ([]Bundle, error)
}
