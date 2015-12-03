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

//Error constants
var (
	ErrBadCRC    = fmt.Errorf("CRC check error")
	ErrWrongSize = fmt.Errorf("size check error")
)

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

func (d *Reader) Read(p []byte) (bytesRead int, err error) {
	defer func() {
		d.Length += bytesRead
	}()

	if d.eof {
		return
	}

	var n int

	if d.buf != nil && len(d.buf) > 0 {
		//Copy and truncate our buffer
		n = copy(p, d.buf)

		if len(d.buf) == n {
			d.buf = nil
		} else {
			d.buf = d.buf[n:]
		}
	} else {
		n, err = d.br.Read(p)
	}

	if err != nil && err != io.EOF {
		return
	}

	lp := len(p)
	var offset int

	//i points at current byte. i-offset is where the current byte should go.
readLoop:
	for i := 0; i < n; i++ {
		switch p[i] {
		case '\r':
			if lp < i+1 {
				return
			}

			if p[i+1] == '\n' {
				//Skip next byte
				i++
				//Set insert position 2 back
				offset += 2
				//Skip this byte
				continue readLoop
			}
		case escape:
			if lp < i+2 {
				return
			}

			if p[i+1] == 'y' {
				var l int
				d.CRC.Write(p[:i])

				l, err = d.checkKeywordLine(p[i:])

				if err != nil {
					if err == io.EOF {
						//Store the remainder of p bytes ( no longer in reader ) into an internal buffer.
						if d.buf != nil && len(d.buf) > 0 {
							d.buf = append(d.buf, p[i:]...)
						} else {
							d.buf = p[i:]
						}
					}
					return
				}

				//Set offset l back
				offset += l
				//Skip all l bytes
				i += l - 1
				continue readLoop
			}

			//Read next byte
			i++
			offset++
			p[i] -= specialOffset
		}

		if !d.begun {
			offset = i + 1
			continue readLoop
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
	preCrc := d.Footer.Get("pcrc32")

	if preCrc == "" {
		return n, nil
	}

	i, err := strconv.ParseUint(preCrc, 16, 0)
	if err != nil {
		return n, fmt.Errorf("error parsing uint: %v", err)
	}

	d.ExpectedCRC = uint32(i)

	length, err := strconv.Atoi(d.Footer.Get("size"))

	if err != nil && length != d.Length {
		return n, ErrWrongSize
	}

	sum := d.CRC.Sum32()
	if sum != uint32(i) {
		return n, ErrBadCRC
	}

	return n, nil
}
