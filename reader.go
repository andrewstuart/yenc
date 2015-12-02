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
	headerLine  = "begin"
	partLine    = "part"
	trailerLine = "end"

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
	begun, eof bool

	Length          int
	CRC             hash.Hash32
	Headers, Footer *YENCHeader
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

	var b byte
	var bs []byte

	var lp = len(p)

readLoop:
	for bytesRead < lp {
		b, err = d.br.ReadByte()

		if err != nil {
			return
		}

		switch b {
		case '\r':
			bs, err = d.br.Peek(1)

			if err != nil {
				return
			}
			if len(bs) < 1 {
				return
			}

			if bs[0] == '\n' {
				_, err = d.br.ReadByte()

				if err != nil {
					return
				}

				continue readLoop
			} else {
				break
			}
		case escape:
			b, err = d.br.ReadByte()

			if err != nil {
				return
			}

			if b == 'y' {
				err = d.checkKeywordLine()
				if err != nil {
					return
				}

				continue readLoop
			}

			b -= specialOffset
		}

		if !d.begun {
			continue readLoop
		}

		p[bytesRead] = b - byteOffset
		d.CRC.Write(p[bytesRead : bytesRead+1])
		bytesRead++
	}

	return
}

func (d *Reader) checkKeywordLine() error {
	bs, err := d.br.Peek(5)

	if err != nil {
		return err
	}

	if beginsWith(bs, headerBytes) || beginsWith(bs, partBytes) {
		d.begun = true

		h, err := ReadYENCHeader(d.br)
		if err != nil {
			return err
		}
		d.Headers = h

		return err
	} else if beginsWith(bs, trailerBytes) {
		d.eof = true

		if err = d.checkTrailer(bs); err != nil {
			return err
		}

		if err == nil {
			err = io.EOF
		}

		return err
	}

	return nil
}

func beginsWith(l, c []byte) bool {
	return len(l) >= len(c) && bytes.Equal(c, l[:len(c)])
}

func (d *Reader) checkTrailer(l []byte) error {
	f, err := ReadYENCHeader(d.br)

	if err != nil {
		return err
	}

	d.Footer = f
	preCrc := d.Footer.Get("pcrc32")

	if preCrc == "" {
		return nil
	}

	i, err := strconv.ParseUint(preCrc, 16, 0)
	if err != nil {
		return fmt.Errorf("error parsing uint: %v", err)
	}

	length, err := strconv.Atoi(d.Footer.Get("size"))

	if err != nil && length != d.Length {
		return ErrWrongSize
	}

	sum := d.CRC.Sum32()
	if sum != uint32(i) {
		return ErrBadCRC
	}

	return nil
}
