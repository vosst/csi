package dmesg

// #include <sys/syslog.h>
import "C"

// Facility models well-known log facilities (see man klogctl)
type Facility int

const (
	LOG_KERN     = C.LOG_KERN     // kernel messages
	LOG_USER     = C.LOG_USER     // random user-level messages
	LOG_MAIL     = C.LOG_MAIL     // mail system
	LOG_DAEMON   = C.LOG_DAEMON   // system daemons
	LOG_AUTH     = C.LOG_AUTH     // security/authorization messages
	LOG_SYSLOG   = C.LOG_SYSLOG   // messages generated internally by syslogd
	LOG_LPR      = C.LOG_LPR      // line printer subsystem
	LOG_NEWS     = C.LOG_NEWS     // network news subsystem
	LOG_UUCP     = C.LOG_UUCP     // UUCP subsystem
	LOG_CRON     = C.LOG_CRON     // clock daemon
	LOG_AUTHPRIV = C.LOG_AUTHPRIV // security/authorization messages (private)
	LOG_FTP      = C.LOG_FTP      // ftp daemon
)

// MaskFacility extracts the facility (bottom 3 bits) from the integer value v.
func MaskFacility(v int) Facility {
	return Facility((v & 0x03fb) >> 3)
}
