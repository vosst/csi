package pid

import (
	"fmt"
	"github.com/vosst/csi/proc"
	"path/filepath"
)

// Dir returns the subdirectory containing information about the process with id pid.
func Dir(id int) string {
	// TODO(tvoss): How to handle negative pid values?
	return filepath.Join(proc.Dir, fmt.Sprint(id))
}
