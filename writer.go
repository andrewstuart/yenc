package yenc

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash"
	"hash/crc32"
	"io"
	"strconv"
)

const (
	null   byte = 0x00
	cr     byte = 0x0a
	lf     byte = 0x0d
	escape byte = 0x3d
	tab    byte = 0x09
	space  byte = 0x20
)

type Writer struct {
	CRC            hash.Hash32
	Length, Line   int
	Name           string
	Header, Footer *YENCHeader

	curLineLen int
	w          io.Writer
	cls        bool
	body, out  *bytes.Buffer
}

func NewWriter(w io.Writer) *Writer {
	b := &bytes.Buffer{}

	return &Writer{
		body: b,
		out:  &bytes.Buffer{},
		w:    w,

		CRC:    crc32.New(crc32.IEEETable),
		Line:   128,
		Header: &YENCHeader{},
		Footer: &YENCHeader{},
	}
}

func (w *Writer) Write(p []byte) (written int, err error) {
	if w.cls {
		return 0, io.ErrClosedPipe
	}

	//Add original to CRC
	w.CRC.Write(p)
	w.Length += len(p)

	for i := range p {
		b := p[i]
		b += byteOffset

		switch b {
		case null, cr, lf, tab, escape, 0x2e:
			w.body.WriteByte(escape)
			b += specialOffset
			w.curLineLen++
		}

		written++
		w.body.WriteByte(b)
		w.curLineLen++

		if w.curLineLen >= w.Line {
			w.body.Write([]byte("\r\n"))
			w.curLineLen = 0
		}
	}

	return
}

func (w *Writer) Close() error {
	w.cls = true

	w.Header.Add("size", strconv.Itoa(w.Length))
	w.Header.Add("name", w.Name)
	w.Header.Add("line", strconv.Itoa(w.Line))

	fmt.Fprintf(w.w, "=ybegin %s \r\n", w.Header)
	_, err := io.Copy(w.w, w.body)

	if err != nil {
		return err
	}

	if w.curLineLen > 0 {
		fmt.Fprint(w.w, "\r\n")
	}

	p := make([]byte, 4)
	binary.BigEndian.PutUint32(p, w.CRC.Sum32())

	w.Footer.Add("pcrc32", hex.EncodeToString(p))
	w.Footer.Add("size", strconv.Itoa(w.Length))

	fmt.Fprintf(w.w, "=yend %s\r\n", w.Footer)

	return nil
}
