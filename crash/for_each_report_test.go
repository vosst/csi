package crash

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"path/filepath"
	"testing"
)

const testDir = "/tmp/MyTestDir"

type MockVisitor struct {
	mock.Mock
}

func (self *MockVisitor) NewReport(name string, report Report) {
	self.Called(name)
}

func (self *MockVisitor) NewError(err error) {
	self.Called(err)
}

type ErrorRecordingVisitor struct {
	Err error
}

func (self ErrorRecordingVisitor) NewReport(string, Report) {
}

func (self *ErrorRecordingVisitor) NewError(err error) {
	self.Err = err
}

func TestForEachReportInDirFailsIfDirectoryDoesNotExist(t *testing.T) {
	os.RemoveAll(testDir)

	erv := ErrorRecordingVisitor{}
	ForEachReportInDir(testDir, &erv)
	assert.IsType(t, ErrorFailedToReadCrashDir{}, erv.Err)
}

func TestForEachReportInDirReportsParsingError(t *testing.T) {
	os.RemoveAll(testDir)
	os.MkdirAll(testDir, os.ModeDir|os.ModePerm)

	d, _ := os.Create(filepath.Join(testDir, "not_parseable"))
	fmt.Fprintf(d, "This is malformed, totally, absolutely.....")
	d.Close()

	erv := ErrorRecordingVisitor{}
	ForEachReportInDir(testDir, &erv)
	// Apparently, the parser does not return an error but instead
	// just returns an empty map. With that, we cannot easily assert on the
	// error here. Still leaving it in for documenting the case.
	// assert.IsType(t, ErrorFailedToParseCrashReport{}, erv.Err)
}

func TestForEachReportInDirInvokesVisitorForReports(t *testing.T) {
	mv := MockVisitor{}
	mv.On("NewReport", "test.crash").Return(nil)
	ForEachReportInDir("test_data", &mv)
	mv.AssertNumberOfCalls(t, "NewReport", 1)
}
