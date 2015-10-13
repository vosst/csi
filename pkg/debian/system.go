package debian

import "github.com/vosst/csi/pkg"

// System implements pkg.System for a Debian system
type System struct {
	dpkg *Dpkg
}

// NewSystem returns a new System instance, providing it with a valid Dpkg instance.
func NewSystem() *System {
	return &System{NewDpkg()}
}

// Resolve returns all packages containing a file matching pattern.
//
// Returns an error if querying the underlying package index fails.
func (self System) Resolve(pattern string) ([]pkg.Bundle, error) {
	packages, err := self.dpkg.QueryForFilePattern(pattern)

	if err != nil {
		return nil, err
	}

	result := []Package{}

	for _, p := range packages {
		if infos, err := self.dpkg.Show(p); err == nil {
			result = append(result, infos)
		}
	}

	bundles := make([]pkg.Bundle, len(result))
	for i, r := range result {
		bundles[i] = r
	}

	return bundles, nil
}

// Arch queries the system architecture that the system has been built for.
//
// Returns an error if querying the information from the system fails.
func (self System) Arch() (pkg.Arch, error) {
	arch, err := self.dpkg.Architecture()
	return pkg.Arch(arch), err
}
