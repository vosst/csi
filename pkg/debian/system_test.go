package debian

import (
	"github.com/stretchr/testify/assert"
	"github.com/vosst/csi/pkg"
	"testing"
)

func TestSystemResolvesCorrectPackages(t *testing.T) {
	system := System{&Dpkg{"test_data"}}

	bundles, err := system.Resolve("/usr/share/*/go*")
	assert.Nil(t, err)
	if assert.Len(t, bundles, 1) {
		assert.Equal(t, "golang", bundles[0].Name())
		assert.Equal(t, "all", bundles[0].Arch())
		assert.Equal(t, "2:1.2.1-2ubuntu1", bundles[0].Version())
	}
}

func TestSystemReturnsCorrectArch(t *testing.T) {
	system := System{&Dpkg{"test_data"}}

	arch, err := system.Arch()
	assert.Nil(t, err)
	assert.Equal(t, pkg.Arch("beautiful"), arch)
}
