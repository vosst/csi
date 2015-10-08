package pid

// #include <elf.h>
import "C"

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
)

const (
	AT_NULL          = C.AT_NULL          // End of vector
	AT_IGNORE        = C.AT_IGNORE        // Entry should be ignored
	AT_EXECFD        = C.AT_EXECFD        // File descriptor of program
	AT_PHDR          = C.AT_PHDR          // Program headers for program
	AT_PHENT         = C.AT_PHENT         // Size of program header entry
	AT_PHNUM         = C.AT_PHNUM         // Number of program headers
	AT_PAGESZ        = C.AT_PAGESZ        // System page size
	AT_BASE          = C.AT_BASE          // Base address of interpreter
	AT_FLAGS         = C.AT_FLAGS         // Flags
	AT_ENTRY         = C.AT_ENTRY         // Entry point of program
	AT_NOTELF        = C.AT_NOTELF        // Program is not elf
	AT_UID           = C.AT_UID           // Real uid
	AT_EUID          = C.AT_EUID          // Effective uid
	AT_GID           = C.AT_GID           // Real gid
	AT_EGID          = C.AT_EGID          // Effective gid
	AT_CLKTCK        = C.AT_CLKTCK        // Frequency of times()
	AT_PLATFORM      = C.AT_PLATFORM      // String identifying platform
	AT_HWCAP         = C.AT_HWCAP         // Machine-dependent hints about processor capabilities
	AT_FPUCW         = C.AT_FPUCW         // Used FPU control word
	AT_DCACHEBSIZE   = C.AT_DCACHEBSIZE   // Data cache block size
	AT_ICACHEBSIZE   = C.AT_ICACHEBSIZE   // Instruction cache block size
	AT_UCACHEBSIZE   = C.AT_UCACHEBSIZE   // Unified cache block size
	AT_IGNOREPPC     = C.AT_IGNOREPPC     // Entry should be ignored
	AT_SECURE        = C.AT_SECURE        // Booleans, was exec setuid-like?
	AT_BASE_PLATFORM = C.AT_BASE_PLATFORM // String identifying real platforms
	AT_RANDOM        = C.AT_RANDOM        // Address of 16 random bytes
	AT_EXECFN        = C.AT_EXECFN        // Filename of executable
	// Pointer to the global system page used for system calls and other nice things
	AT_SYSINFO      = C.AT_SYSINFO
	AT_SYSINFO_EHDR = C.AT_SYSINFO_EHDR
	// Shapes of the caches, bits 0-3 contains associativity; bits 4-7 contains log2 of line size. Maks those to get cache size
	AT_L1I_CACHESHAPE = C.AT_L1I_CACHESHAPE
	AT_L1D_CACHESHAPE = C.AT_L1D_CACHESHAPE
	AT_L2_CACHESHAPE  = C.AT_L2_CACHESHAPE
	AT_L3_CACHESHAPE  = C.AT_L3_CACHESHAPE
)

func determineEndianess() binary.ByteOrder {
	var x uint32 = 0x01020304

	switch *(*byte)(unsafe.Pointer(&x)) {
	case 0x01:
		return binary.BigEndian
	case 0x04:
		return binary.LittleEndian
	}

	return nil
}

// Describes the contents of the ELF interpreter information passed to the process at exec time.
type Auxv map[uint32]uint32

// NewAuxv reads the /proc/pid/auxv entry and returns the correspoding Auxv instance if
// reading the file was successful or an error otherwise.
func NewAuxv(pid int) (Auxv, error) {
	fn := filepath.Join(Dir(pid), "auxv")

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

	return NewAuxvFromReader(bytes.NewReader(b)), nil
}

// NewAuxvFromReader reads the contents from reader returning
// a corresponding Auxv instance or nil in case of issues.
func NewAuxvFromReader(reader io.Reader) Auxv {
	endianess := determineEndianess()

	if endianess == nil {
		return nil
	}

	auxv := Auxv(make(map[uint32]uint32))

	pair := []uint32{0, 0}

	for err := binary.Read(reader, endianess, &pair); err == nil; err = binary.Read(reader, endianess, &pair) {
		if pair[0] == AT_NULL {
			break
		}

		auxv[pair[0]] = pair[1]
	}

	return auxv
}
