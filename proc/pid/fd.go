package pid

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// File describes an fd referring to an actual file
type File string

// SocketOrPipe describes an fd refering to a socket or pipe
type SocketOrPipe struct {
	Type  string // The type, either socket or pipe
	Inode uint   // The INode
}

// Anon describes an anonymous fd without an inode
type Anon string

// Fd is a slice of elements each representing an fd being in use by a process.
//
// Depending on the type of an individual fd, the element is either a File, a SocketOrPipe or
// an Anon instance.
type Fd []interface{}

// NewFd returns all open fds for the process identified by pid.
//
// Returns an error if opening /proc/%{pid}/fd or a subsequent os.File.Readdir failed
func NewFd(pid int) (Fd, error) {
	fn := filepath.Join(Dir(pid), "fd")

	f, err := os.Open(fn)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to open %s [%s]", fn, err))
	}

	defer f.Close()

	if entries, err := f.Readdir(0); err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to read %s [%s]", fn, err))
	} else {
		fd := Fd{}

		for _, fi := range entries {
			fni := filepath.Join(fn, fi.Name())
			dest, _ := os.Readlink(fni)
			// We could not stat the destination and assume that the fd is not
			// an actual file but either a socket/pipe or an fd without an inode.
			if _, err := os.Stat(dest); err != nil {
				kv := strings.Split(dest, ":")

				if len(kv) != 2 {
					continue
				}

				// We got an fd without an inode, parsing the file type now.
				if kv[0] == "anon_inode" {
					fd = append(fd, Anon(strings.Trim(kv[1], "[]")))
				} else {
					inode := uint(0)
					fmt.Sscanf(kv[1], "[%d]", &inode)
					fd = append(fd, SocketOrPipe{kv[0], inode})
				}
			} else {
				fd = append(fd, File(dest))
			}
		}

		return fd, nil
	}
}
