package pid

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Pathname of the executed command
type Exe string

// NewExe determines the path of the executed command from /proc/pid/exe.
// Returns an error if following the symbolic link fails.
func NewExe(pid int) (Exe, error) {
	fn := filepath.Join(Dir(pid), "exe")

	lr, err := os.Readlink(fn)
	if err != nil {
		return Exe(lr), errors.New(fmt.Sprintf("Failed to resolve link %s to cwd [%s].", fn, err))
	}

	return Exe(lr), nil

}
