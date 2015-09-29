package main

import "io"

// NewLineReader reads a final \n if the contained Reader
// reports EOF.
type NewLineReader struct {
	Reader io.Reader
}

// Read triggers the internal Reader instance. When that instance
// reports EOF, an additional \n is reported back to the caller.
func (self NewLineReader) Read(p []byte) (int, error) {
	n, err := self.Reader.Read(p)
	if err == io.EOF {
		p[n] = '\n'
		return n + 1, err
	}
	return n, err
}
