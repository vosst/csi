package crash

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/vosst/csi/machine"
	"labix.org/v2/mgo/bson"
)

// We have to submit a whoopsie version field when uploading to
// the crash handling infrastructure.
const whoopsieVersion = "0.2.49"

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

// HttpReporterPersister persists incoming crash reports to launchpad.
type HttpReportPersister struct {
	SubmitURL  url.URL            // URL for sending the crash report to
	Identifier machine.Identifier // Identifier helps in generating a globally unique device id
	Client     *http.Client       // HTTP client instance for reaching out to the crash db service
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

func (self HttpReportPersister) uploadCore(report Report, oopsId string) error {
	core, contains := report["CoreDump"]
	if !contains || len(core) == 0 {
		return errors.New("Missing field CoreDump in report")
	}

	arch, contains := report["Architecture"]
	if !contains || len(arch) == 0 {
		return errors.New("Missing field Architecture in report")
	}

	id, err := self.Identifier.Identify()

	if err != nil {
		return err
	}

	coreURL := fmt.Sprintf("%s/%s/submit-core/%s/%s", self.SubmitURL, oopsId, arch, hex.EncodeToString(id))
	req, err := http.NewRequest("POST", coreURL, strings.NewReader(core[0]))
	if err != nil {
		return err
	}

	// TODO(vosst): add handling of response here.
	_, err = self.Client.Do(req)
	return err
}

func (self HttpReportPersister) handleUploadResponse(report Report, response string) error {
	if len(response) == 0 {
		return nil
	}

	var oopsId, command string
	if _, err := fmt.Sscanf(response, "%s %s", &oopsId, &command); err != nil {
		return errors.New("Failed to parse response body")
	}

	switch command {
	case "CORE":
		return self.uploadCore(report, oopsId)
	case "OOPSID":
		log.Printf("Server reported OOPS ID: %s", oopsId)
	}

	return nil
}

func (self HttpReportPersister) Persist(report Report) error {
	bson, _ := self.marshalToBSON(report)

	if resp, err := self.Client.Post(self.SubmitURL.String(), "application/octet-stream", bytes.NewReader(bson)); err != nil {
		fmt.Print(err)
		return err
	} else if resp.StatusCode == http.StatusOK {
		// We received a response and try to interpret it further
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err != nil {
			return err
		} else {
			return self.handleUploadResponse(report, string(body))
		}
	} else {
		return errors.New(fmt.Sprintf("Received status code %d, indicating an issue with our upload", resp.StatusCode))
	}
}
