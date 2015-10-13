package debian

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPackageQueriesPackageKeyForName(t *testing.T) {
	pkg := Package{}
	pkg["Package"] = []string{"lalelu"}
	assert.Equal(t, "lalelu", pkg.Name())
}

func TestPackageQueriesArchitectureKeyForArch(t *testing.T) {
	pkg := Package{}
	pkg["Architecture"] = []string{"beautiful"}
	assert.Equal(t, "beautiful", pkg.Arch())
}

func TestPackageQueriesVersionKeyForVersion(t *testing.T) {
	pkg := Package{}
	pkg["Version"] = []string{"1.2.3"}
	assert.Equal(t, "1.2.3", pkg.Version())
}
