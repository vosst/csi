package pid

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

type Hierarchy struct {
	ID           int      // Hierarchy ID number
	Subsystems   []string // Set of subsystems bound to the hierarchy
	ControlGroup string   // Control group in the hierarchy to which the process belongs
}

func NewHierarchyFromLine(line string) (Hierarchy, error) {
	h := Hierarchy{}
	s := ""

	if _, err := fmt.Sscanf(line, "%d:%s:%s", &h.ID, &s, &h.ControlGroup); err != nil {
		return h, errors.New(fmt.Sprintf("Failed to parse line '%s' [%s]", line, err))
	}

	h.Subsystems = strings.Split(s, ",")
	return h, nil
}

// Cgroup describes control group Hierarchies to which the process/task belongs.
type Cgroup struct {
	Hierarchies []Hierarchy
}

func NewCGroup(pid int) (*Cgroup, error) {
	fn := filepath.Join(Dir(pid), "cgroup")

	f, err := os.Open(fn)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to read %s [%s]", fn, err))
	}

	defer f.Close()

	stat, err := f.Stat()

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to stat %s [%s]", fn, err))
	}

	b, err := syscall.Mmap(int(f.Fd()), 0, int(stat.Size()), syscall.PROT_READ, syscall.MAP_PRIVATE)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to mmap %s [%s]", fn, err))
	}

	defer syscall.Munmap(b)

	return NewCgroupFromReader(bytes.NewReader(b)), nil
}

func NewCgroupFromReader(reader io.Reader) *Cgroup {
	cg := Cgroup{}
	br := bufio.NewReader(reader)

	for line, err := br.ReadString('\n'); err == nil; line, err = br.ReadString('\n') {
		if h, err := NewHierarchyFromLine(line); err == nil {
			cg.Hierarchies = append(cg.Hierarchies, h)
		}
	}

	return &cg
}
