package main

import (
	"bufio"
	"io"
	"net/textproto"
)

// CrashReport models an apport crash report.
type CrashReport textproto.MIMEHeader

// ParseCrashReport parses information from reader.
func ParseCrashReport(reader io.Reader) (CrashReport, error) {
	tr := textproto.NewReader(bufio.NewReader(reader))

	hdr, err := tr.ReadMIMEHeader()

	if err != nil {
		return nil, err
	}

	return CrashReport(hdr), nil
}
