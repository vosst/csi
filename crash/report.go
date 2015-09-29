package crash

import (
	"bufio"
	"io"
	"net/textproto"
)

// Report models an apport crash report.
type Report textproto.MIMEHeader

// ParseReport parses information from reader.
func ParseReport(reader io.Reader) (Report, error) {
	tr := textproto.NewReader(bufio.NewReader(reader))

	hdr, err := tr.ReadMIMEHeader()

	if err != nil {
		return nil, err
	}

	return Report(hdr), nil
}
