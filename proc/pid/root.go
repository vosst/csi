package pid

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Root of the filesystem as seen by a process
type Root string

// NewRoot determines the filsystem root of the process identified by pid
//
// Returns an error if resolving the symbolic link fails, usually caused by
// the process's main thread having exited already.
func NewRoot(pid int) (Root, error) {
	fn := filepath.Join(Dir(pid), "root")

	if r, err := os.Readlink(fn); err != nil {
		return "", errors.New(fmt.Sprintf("Failed to resolve symbolic link [%s]", err))
	} else {
		return Root(r), nil
	}
}
