package crash

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// ErrorFailedToReadCrashDir indicates an issue when trying to read entries of
// the given crash dir.
type ErrorFailedToReadCrashDir struct {
	CrashDir string // The crash dir that we failed to read.
	Next     error  // The error reported when trying to read the directory.
}

// Error pretty prints the given ErrorFailedToReadCrashDir instance.
func (self ErrorFailedToReadCrashDir) Error() string {
	return fmt.Sprintf("Failed to read crash dir %s [%s]", self.CrashDir, self.Next)
}

// ErrorFailedToParseCrashReport indicates an issue trying to parse a crash report.
type ErrorFailedToParseCrashReport struct {
	Name string // Name of the report that failed parsing
	Next error  // Error reported by the parser.
}

// Error pretty prints the given ErrorFailedToParseCrashReport instance.
func (self ErrorFailedToParseCrashReport) Error() string {
	return fmt.Sprintf("Failed to parse crash report %s [%s]", self.Name, self.Next)
}

// ReportVisitor abstracts handling of individual Report instances.
type ReportVisitor interface {
	// NewReport is invoked for every report that has been successfully loaded
	NewReport(name string, report Report)
	// NewError is invoked for issues encountered while processing individual entries.
	NewError(err error)
}

// ForEachReportInDir iterates over all crash reports in dir, parses it and reports
// back to visitor.
func ForEachReportInDir(dir string, visitor ReportVisitor) {
	entries, err := ioutil.ReadDir(dir)

	if err != nil {
		visitor.NewError(ErrorFailedToReadCrashDir{dir, err})
		return
	}

	for _, entry := range entries {
		fn := filepath.Join(dir, entry.Name())

		if matched, _ := filepath.Match("*.crash", entry.Name()); !matched {
			// Silently skip over the file.
			continue
		}

		f, err := os.Open(fn)
		if err != nil {
			visitor.NewError(ErrorFailedToParseCrashReport{entry.Name(), err})
			continue
		}

		defer f.Close()
		report, err := ParseReport(NewLineReader{f})

		if err != nil {
			visitor.NewError(ErrorFailedToParseCrashReport{entry.Name(), err})
			continue
		}

		visitor.NewReport(entry.Name(), report)
	}
}
