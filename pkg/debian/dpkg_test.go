package debian

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDpkgFindsPackagesMatchingPatterns(t *testing.T) {
	dpkg := Dpkg{"test_data"}

	// We expect exactly one package to be found here.
	pkgs, err := dpkg.QueryForFilePattern("/usr/share/doc/go*/*")
	assert.Nil(t, err)
	assert.Equal(t, "golang", pkgs[0])

	// And no package for me, myself and I.
	pkgs, err = dpkg.QueryForFilePattern("/me/myself/and/I")
	assert.Nil(t, err)
	assert.Len(t, pkgs, 0)
}

func TestDpkgShowsPackagesCorrectly(t *testing.T) {
	dpkg := Dpkg{"test_data"}
	pkg, err := dpkg.Show("golang")

	assert.Nil(t, err)
	if assert.NotNil(t, pkg) {
		assert.Equal(t, true, pkg.IsInstalledCorrectly())
		assert.Equal(t, "golang", pkg.Name())
		assert.Equal(t, "2:1.2.1-2ubuntu1", pkg.Version())
		assert.Equal(t, "all", pkg.Arch())
	}
}
