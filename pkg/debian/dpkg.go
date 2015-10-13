package debian

import (
	"bufio"
	"errors"
	"fmt"
	"net/textproto"
	"os"
	"path/filepath"
	"regexp"
)

// extractPnRegexp extracts the package name.
//
// Example:
//   /var/lib/dpkg/info/some-package:amd64.list
//                      ^^^^^^^^^^^^
// extracts the marked part, with submatch index 1.
var extractPnRegexp = regexp.MustCompile(`.*\/([^\:]+)(\:[^\:]+)?\.list`)

// Submatch index of the package name in extractPnRegexp
const smIdxPn = 1

// Dpkg provides access to a subset of the overall dpkg features.
//
// Instead of exec'ing dpkg-query (which is sloooow), we instead rely on
// a native implementation that allows us to execute the operations we need
// fast and asynchronously.
type Dpkg struct {
	runtimeDir string // Runtime directory containing dpkg's files
}

// NewDpkg returns a new Dpkg instance, pointing to the system default dpkg runtime dir
func NewDpkg() *Dpkg {
	return &Dpkg{"/var/lib/dpkg"}
}

func (self Dpkg) Architecture() (string, error) {
	archFn := filepath.Join(self.runtimeDir, "arch")

	f, err := os.Open(archFn)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to open arch file %s [%s]", archFn, err))
	}

	tpr := textproto.NewReader(bufio.NewReader(f))

	line, err := tpr.ReadLine()
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to read first line of arch file %s [%s]", archFn, err))
	}

	return line, nil
}

// QueryForFilePattern searches through all installed files, matching them against
// pattern and returns the list of package names containing a file matching patterns.
//
// Returns an error if globbing /var/lib/dpkg/info/*.list fails.
func (self Dpkg) QueryForFilePattern(pattern string) ([]string, error) {
	result := []string{}

	info := filepath.Join(self.runtimeDir, "info")

	entries, err := filepath.Glob(filepath.Join(info, "*.list"))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to glob *.list from %s [%s]", info, err))
	}

	for _, entry := range entries {
		if f, err := os.Open(entry); err != nil {
			continue
		} else {
			tpr := textproto.NewReader(bufio.NewReader(f))

			for line, err := tpr.ReadLine(); err == nil; line, err = tpr.ReadLine() {
				if matched, _ := filepath.Match(pattern, line); matched {
					if matches := extractPnRegexp.FindStringSubmatch(entry); len(matches) > 0 {
						result = append(result, matches[smIdxPn])
						// We have to break the inner loop as we want to avoid
						// adding the same package over and over again.
						break
					}
				}
			}
		}
	}

	return result, nil
}

// Show loads all package information for the package with name.
func (self Dpkg) Show(name string) (Package, error) {
	status := filepath.Join(self.runtimeDir, "status")

	f, err := os.Open(status)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to open status file %s [%s]", status, err))
	}

	bf := bufio.NewReader(f)
	for pkg, err := NewPackage(bf); err == nil; pkg, err = NewPackage(bf) {
		if pkg.IsInstalledCorrectly() && pkg.Name() == name {
			return pkg, nil
		}
	}

	return Package{}, nil
}
