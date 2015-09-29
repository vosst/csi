package crash

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"labix.org/v2/mgo/bson"
)

// All the fields we will unconditionally accept for upload,
// no matter their size.
var AcceptableFields = map[string]struct{}{
	"ProblemType":                struct{}{},
	"Date":                       struct{}{},
	"Traceback":                  struct{}{},
	"Signal":                     struct{}{},
	"PythonArgs":                 struct{}{},
	"Package":                    struct{}{},
	"SourcePackage":              struct{}{},
	"PackageArchitecture":        struct{}{},
	"Dependencies":               struct{}{},
	"MachineType":                struct{}{},
	"StacktraceAddressSignature": struct{}{},
	"ApportVersion":              struct{}{},
	"DuplicateSignature":         struct{}{},
	"DistroRelease":              struct{}{},
	"Uname":                      struct{}{},
	"Architecture":               struct{}{},
	"NonfreeKernelModules":       struct{}{},
	"LiveMediaBuild":             struct{}{},
	"UserGroups":                 struct{}{},
	"ExecutablePath":             struct{}{},
	"InterpreterPath":            struct{}{},
	"ExecutableTimestamp":        struct{}{},
	"ProcCwd":                    struct{}{},
	"ProcEnviron":                struct{}{},
	"ProcCmdline":                struct{}{},
	"ProcStatus":                 struct{}{},
	"ProcMaps":                   struct{}{},
	"ProcAttrCurrent":            struct{}{},
	"Registers":                  struct{}{},
	"Disassembly":                struct{}{},
	"StacktraceTop":              struct{}{},
	"AssertionMessage":           struct{}{},
	"CoreDump":                   struct{}{},
	"VmCore":                     struct{}{},
	"Tags":                       struct{}{},
	"OopsText":                   struct{}{},
	"UpgradeStatus":              struct{}{},
	"InstallationDate":           struct{}{},
	"InstallationMedia":          struct{}{},
	"IwlFwDump":                  struct{}{},
	"SystemImageInfo":            struct{}{},
}

// We always filter out these fields and do not submit them
// to the cloud.
var UnacceptableFields = map[string]struct{}{
	"Stacktrace":         struct{}{}, // Not needed, we retrace with ddebs on errors.
	"ThreadStacktrace":   struct{}{}, // Not needed, we retrace with ddebs on errors.
	"UnreportableReason": struct{}{}, // No use for it right now.
	"CrashCounter":       struct{}{}, // We maintain our own count.
	"_MarkForUpload":     struct{}{}, // Redundant since the crash was uploaded.
	"Title":              struct{}{},
}

var ErrHostNotReachable = errors.New("Destination is not reachable")
var ErrWillNotSendViaWWAN = errors.New("Destination Host is only available via WWAN.")

// HttpReporterPersister persists incoming crash reports to launchpad.
type HttpReportPersister struct {
	SubmitURL       url.URL // URL for sending the crash report to
	WhoopsieVersion string  // WhoopsieVersion field that is communicated to the server on upload.
	// ReachabilityMonitor ReachabilityMonitor // Monitors whether SubmitURL is reachable.
	Client *http.Client // HTTP client instance for reaching out to the crash db service
}

// filterField returns true if the given (key, value) pair should be filtered out.
func (self HttpReportPersister) filterField(k string, v []string) bool {
	if len(v) == 0 {
		return true
	}

	if _, present := UnacceptableFields[k]; present {
		return true
	}

	if _, present := AcceptableFields[k]; present {
		return false
	}

	return len(v[0]) > 1024*1024
}

// marshalToBSON walks the given report, filtering out all invalid fields
// and encodes the resulting filtered map to BSON.
func (self HttpReportPersister) marshalToBSON(report Report) ([]byte, error) {
	filtered := make(map[string]string)

	for k, v := range report {
		if !self.filterField(k, v) && len(v) > 0 {
			filtered[k] = v[0]
		}
	}

	return bson.Marshal(filtered)
}

func (self HttpReportPersister) Persist(report Report) error {
	/*reachability := self.ReachabilityMonitor.CheckHostReachability(self.SubmitURL.Host)

	if reachability == NotReachable {
		return ErrHostNotReachable
	}

	if reachability&IsWWAN == IsWWAN {
		return ErrWillNotSendViaWWAN
	}
	*/

	bson, _ := self.marshalToBSON(report)
	req, _ := http.NewRequest("POST", self.SubmitURL.String(), bytes.NewReader(bson))
	req.Header.Add("X-Whoopsie-Version", self.WhoopsieVersion)

	if resp, err := self.Client.Do(req); err == nil {
		// We received a response and try to interpret it further
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		fields := strings.Split(string(body), " ")

		switch fields[0] {
		case "CORE":
			log.Print("Server requested upload of coredump")
			core := report["CoreDump"][0]
			arch := report["Architecture"][0]

			if len(core) > 0 && len(arch) > 0 {
				coreURL := fmt.Sprintf("%s/%s/submit-core/%s/%s", self.SubmitURL, "uuid", arch, "id")
				req, err = http.NewRequest("POST", coreURL, strings.NewReader(core))
				req.Header.Add("X-Whoopsie-Version", self.WhoopsieVersion)
				resp, err = self.Client.Do(req)

			}
		case "OOPSID":
			log.Printf("Server reported OOPS ID: %s", fields[1])
		default:
			log.Printf("Received unhandled command: %s", string(body))
		}
	}

	return nil
}
