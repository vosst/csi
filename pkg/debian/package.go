package debian

import (
	"bufio"
	"errors"
	"fmt"
	"net/textproto"
)

type Package textproto.MIMEHeader

func NewPackage(reader *bufio.Reader) (Package, error) {
	tpr := textproto.NewReader(reader)

	hdr, err := tpr.ReadMIMEHeader()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to parse package information [%s]", err))
	}

	return Package(hdr), nil
}

func (self Package) Name() string {
	if v, present := self["Package"]; present && len(v) > 0 {
		return v[0]
	}

	return ""
}

func (self Package) Version() string {
	if v, present := self["Version"]; present && len(v) > 0 {
		return v[0]
	}

	return ""
}

func (self Package) Arch() string {
	if v, present := self["Architecture"]; present && len(v) > 0 {
		return v[0]
	}

	return ""
}

func (self Package) IsInstalledCorrectly() bool {
	if v, present := self["Status"]; present && len(v) > 0 {
		return v[0] == "install ok installed"
	}

	return false
}
