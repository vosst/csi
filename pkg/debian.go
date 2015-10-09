package pkg

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net/textproto"
	"os/exec"
	"strings"
)

// DebInfo is a dictionary describing a debian package
type DebInfo map[string]string

// Name returns the binary package name
func (self DebInfo) Name() string {
	return self["binary:Package"]
}

// Version returns the package version
func (self DebInfo) Version() string {
	return self["Version"]
}

// Arch returns the architecture that the binary package has been built for
func (self DebInfo) Arch() string {
	return self["Architecture"]
}

// Dpkg executes dpkg
type Dpkg struct {
}

// PrintArchitectures returns the system's architecture.
//
// Returns an error if executing dpkg fails.
func (self Dpkg) PrintArchitecture() (Arch, error) {
	out, err := exec.Command("dpkg", "--print-architecture").Output()

	if err != nil {
		return Arch(""), errors.New(fmt.Sprintf("Failed to query system architecture [%s]", err))
	}

	return Arch(strings.TrimSpace(string(out))), err
}

// defaultFields lists the fields we would like to include in reports.
// Field names are taken from man dpkg-query.
// Please make sure to adhere to dpkg-query's format specification when adding to this list.
var defaultFields = []string{
	"${binary:Package}",
	"${source:Package}",
	"${Architecture}",
	"${Pre-Depends}",
	"${Depends}",
	"${Origin}",
	"${Version}",
}

// DpkgQuery executes dpkg-query.
type DpkgQuery struct {
}

// Search queries the package index with pattern, returning the names of all packages that contain a file matching the pattern.
//
// Returns an error if executing dpkg-query fails.
func (self DpkgQuery) Search(pattern string) ([]string, error) {
	out, err := exec.Command("dpkg-query", "-S", pattern).Output()

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to search the package index [%s]", err))
	}

	pkgs := []string{}
	tr := textproto.NewReader(bufio.NewReader(bytes.NewReader(out)))

	for line, err := tr.ReadLine(); err == nil; line, err = tr.ReadLine() {
		if kv := strings.Split(line, ":"); len(kv) == 2 {
			pkgs = append(pkgs, kv[0])
		}
	}

	return pkgs, nil
}

// Show queries the packaging index for all packages matching the given patterns.
// Returns a slice of DebianPackageInfo instances, with each instance containing the subset of
// fields returned by the packaging index.
//
// Returns an error if querying the packaging index fails.
func (self DpkgQuery) Show(pattern string, fields []string) ([]DebInfo, error) {
	format := strings.Join(fields, "|")

	out, err := exec.Command("dpkg-query", "-f", format, "-W", pattern).Output()

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to show information for package %s [%s]", pattern, err))
	}

	infos := []DebInfo{}
	tr := textproto.NewReader(bufio.NewReader(bytes.NewReader(out)))

	for line, err := tr.ReadLine(); err == nil; line, err = tr.ReadLine() {
		pi := DebInfo{}

		for i, v := range strings.Split(line, "|") {
			if len(v) > 0 {
				pi[fields[i]] = v
			}
		}

		infos = append(infos, pi)
	}

	return infos, nil
}

// DebianSystem implements pkg.System for a Debian system
type DebianSystem struct {
}

// Resolve returns all packages containing a file matching pattern.
//
// Returns an error if querying the underlying package index fails.
func (self DebianSystem) Resolve(pattern string) ([]Bundle, error) {
	dpkgQuery := DpkgQuery{}
	packages, err := dpkgQuery.Search(pattern)

	if err != nil {
		return nil, err
	}

	result := []DebInfo{}

	for _, p := range packages {
		if infos, err := dpkgQuery.Show(p, defaultFields); err == nil {
			result = append(result, infos...)
		}
	}

	bundles := make([]Bundle, len(result))
	for i, r := range result {
		bundles[i] = r
	}

	return bundles, nil
}

// Arch queries the system architecture that the system has been built for.
//
// Returns an error if querying the information from the system fails.
func (self DebianSystem) Arch() (Arch, error) {
	dpkg := Dpkg{}
	return dpkg.PrintArchitecture()
}
