package pid

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Cwd describes the current working directory of the process.
type Cwd string

// NewCwd determines the current working directory of the process identified by pid,
// returning an error if following the link to the cwd fails.
func NewCwd(pid int) (Cwd, error) {
	fn := filepath.Join(Dir(pid), "cwd")

	lr, err := os.Readlink(fn)
	if err != nil {
		return Cwd(lr), errors.New(fmt.Sprintf("Failed to resolve link %s to cwd [%s].", fn, err))
	}

	return Cwd(lr), nil
}
