package yenc

import (
	"bufio"
	"bytes"
	"fmt"
	"hash"
	"hash/crc32"
	"io"
	"strconv"
)

const (
	keywordLine = "=y"
	headerLine  = "=ybegin"
	partLine    = "=ypart"
	trailerLine = "=yend"

	byteOffset    byte = 42
	specialOffset byte = 64
)

var (
	keywordBytes = []byte(keywordLine)
	headerBytes  = []byte(headerLine)
	partBytes    = []byte(partLine)
	trailerBytes = []byte(trailerLine)
)

//ErrWrongSize is returned when a part or file is not the expected full size
type ErrWrongSize struct {
	Expected, Actual int
}

// Error implements error
func (e ErrWrongSize) Error() string {
	return fmt.Sprintf("size check error, expected %d, actual %d", e.Expected, e.Actual)
}

// ErrBadCRC is returned when a CRC does not match
type ErrBadCRC struct {
	Expected, Actual uint32
}

// Error implements error
func (e ErrBadCRC) Error() string {
	return fmt.Sprintf("bad crc; expected %d got %d", e.Expected, e.Actual)
}

//Reader implements the io.Reader methods for an underlying YENC
//document/stream. It additionally exposes some of the metadata that may be
//useful for consumers.
type Reader struct {
	br         *bufio.Reader
	buf        []byte
	begun, eof bool

	Length          int
	ExpectedCRC     uint32
	CRC             hash.Hash32
	Headers, Footer *Header
}

//NewReader returns a reader from an input reader.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		br:  bufio.NewReader(r),
		CRC: crc32.New(crc32.IEEETable),
	}
}

func (d *Reader) buffer(bs []byte) {
	//Store the remainder of p bytes ( no longer in reader ) into an internal buffer.
	if d.buf != nil && bs != nil && len(d.buf) > 0 {
		d.buf = append(d.buf, bs...)
	} else {
		d.buf = bs
	}
}

func (d *Reader) Read(p []byte) (bytesRead int, err error) {
	defer func() {
		d.Length += bytesRead
	}()

	if d.eof {
		err = io.EOF
		return
	}

	n, m := 0, 0
	if d.buf != nil && len(d.buf) > 0 {
		//Copy and truncate our buffer
		n = copy(p, d.buf)
		d.buffer(nil)
	}

	m, err = d.br.Read(p[n:])
	n += m

	if err != nil && err != io.EOF {
		return
	}

	var offset int

	//i points at current byte. i-offset is where the current byte should go.
	for i := 0; i < n; i++ {
		switch p[i] {
		case '\r':
			// If i+1 is further than we've read
			if n < i+1 {
				d.CRC.Write(p[:bytesRead])
				return
			}

			if i+1 > len(p)-1 {
				d.CRC.Write(p[:bytesRead])
				d.buffer(p[i:n])
				return
			}

			if p[i+1] == '\n' {
				//Skip next byte
				i++
				//Set insert position 2 back
				offset += 2
				//Skip this byte
				continue
			}
		case escape:
			if n < i+2 {
				d.CRC.Write(p[:bytesRead])
				d.buffer(p[i:n])
				return
			}

			if p[i+1] == 'y' {
				var l int

				// What if keywordLine is across a boundary?
				if !bytes.Contains(p[i:n], []byte("\r\n")) {
					d.CRC.Write(p[:bytesRead])

					d.buffer(p[i:n])
					return
				}

				l, err = d.checkKeywordLine(p[i:n])

				if err != nil {
					d.CRC.Write(p[:bytesRead])

					if err == io.EOF {
						d.eof = true
						d.buffer(p[i:n])

						err = d.CheckCRC()
					}
					return
				}

				//Set offset l back
				offset += l
				//Skip all l bytes to ignore header line
				i += l - 1
				continue
			}

			//Read next byte
			i++
			offset++
			p[i] -= specialOffset
		}

		if !d.begun {
			offset = i + 1
			continue
		}

		p[i-offset] = p[i] - byteOffset
		bytesRead++
	}

	d.CRC.Write(p[:bytesRead])
	return
}

func (d *Reader) checkKeywordLine(bs []byte) (n int, err error) {
	if bytes.HasPrefix(bs, headerBytes) || bytes.HasPrefix(bs, partBytes) {
		d.begun = true
		d.CRC.Reset()

		var h *Header
		h, n = ReadYENCHeader(bs)
		d.Headers = h
		return
	}

	if bytes.HasPrefix(bs, trailerBytes) {
		d.eof = true

		if n, err = d.checkTrailer(bs); err != nil {
			return
		}

		if err == nil {
			err = io.EOF
		}
	}

	return
}

func (d *Reader) checkTrailer(l []byte) (int, error) {
	f, n := ReadYENCHeader(l)
	d.Footer = f

	length, err := strconv.Atoi(d.Footer.Get("size"))

	if err != nil && length != d.Length {
		return n, ErrWrongSize{Expected: length, Actual: d.Length}
	}

	preCrc := d.Footer.Get("pcrc32")
	if preCrc == "" {
		return n, nil
	}

	i, err := strconv.ParseUint(preCrc, 16, 0)
	if err != nil {
		return n, fmt.Errorf("error parsing uint: %v", err)
	}

	d.ExpectedCRC = uint32(i)
	return n, nil
}

// CheckCRC makes sure the CRC matches and returns an error if not
func (d *Reader) CheckCRC() error {
	sum := d.CRC.Sum32()

	if sum != d.ExpectedCRC {
		return &ErrBadCRC{Actual: sum, Expected: d.ExpectedCRC}
	}
	return nil
}
