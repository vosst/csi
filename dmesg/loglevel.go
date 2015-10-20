package dmesg

// #include <sys/syslog.h>
import "C"

// Loglevel models the kernel's loglevel (see man klogctl)
type Loglevel uint

const (
	LOG_EMERG   Loglevel = C.LOG_EMERG   // system is unusable
	LOG_ALERT            = C.LOG_ALERT   // action must be taken immediately
	LOG_CRIT             = C.LOG_CRIT    // critical conditions
	LOG_ERR              = C.LOG_ERR     // error conditions
	LOG_WARNING          = C.LOG_WARNING // warning conditions
	LOG_NOTICE           = C.LOG_NOTICE  // normal but significant condition
	LOG_INFO             = C.LOG_INFO    // informational
	LOG_DEBUG            = C.LOG_DEBUG   // debug-level-messages
)

// MaskLoglevel extracts the Loglevel from the integer value v
func MaskLoglevel(v uint) Loglevel {
	return Loglevel(v & 0x07)
}
