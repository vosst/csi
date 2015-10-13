package pkg

// Arch describes a machine architecture
type Arch string

const (
	ArchAmd64  Arch = "amd64"  // 64-bit PC
	ArchArmel       = "armel"  // EABI ARM
	ArchArmhf       = "armhf"  // Hard float ABI ARM
	ArchI386        = "i386"   // 32-bit PC
	ArchIa64        = "ia64"   // Intel Itanium
	ArchMips        = "mips"   // MIPS (big-endian mode)
	ArchMipsel      = "mipsel" // MIPS (little-endian mode)
	ArchPPC         = "ppc"    // Motorola/IBM PowerPC
	ArchPPC64       = "ppc64"  // POWER7+, POWER8
)
