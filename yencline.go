package yenc

import (
	"bytes"
	"fmt"
	"strings"
)

//Header is a map type for reading a serialized yenc header.
type Header map[string]string

//String implements the stringer interface, printing the proper yenc header
func (y *Header) String() string {
	s := make([]string, 0, 5)

	for k, v := range *y {
		s = append(s, fmt.Sprintf("%s=%s", k, v))
	}

	return strings.Join(s, " ")
}

//Put creates a string entry for a k, v pair
func (y *Header) Put(k, v string) {
	(*y)[k] = v
}

//Get returns the string for a key
func (y *Header) Get(k string) string {
	return (*y)[k]
}

//ReadYENCHeader accepts a byte slice and returns a YENCHeader or any error
//encountered while decoding, and the header length so that the consumer can
//ignore the appropriate bytes.
func ReadYENCHeader(bs []byte) (*Header, int) {
	i := bytes.IndexByte(bs, '\n') + 1
	s := string((bs)[:i])

	s = strings.TrimSpace(s)

	y := &Header{}

	ss := strings.Split(s, " ")
	for _, kvString := range ss {
		kvPair := strings.Split(kvString, "=")

		if len(kvPair) == 2 {
			y.Put(kvPair[0], kvPair[1])
		}
	}

	return y, i
}
