package yenc

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/andrewstuart/nntp"
)

func TestDecMulti(t *testing.T) {
	b := &bytes.Buffer{}

	f1, err := os.Open("./test/00000020.ntx")

	if err != nil {
		t.Fatalf("f1 error")
	}

	art, err := nntp.NewResponse(f1)

	if err != nil {
		t.Errorf("Article 1 err: %v\n", err)
	}

	r := NewReader(art.Body)

	_, err = bufio.NewReader(r).WriteTo(b)
	if err != nil {
		if err == ErrBadCRC {
			t.Logf("expected crc: %d, was %d.\n", r.ExpectedCRC, r.CRC.Sum32())
		}
		t.Errorf("Article 1 error: %v\n", err)
	}

	f2, err := os.Open("./test/00000021.ntx")

	if err != nil {
		t.Fatalf("f2 error")
	}

	art2, err := nntp.NewResponse(f2)

	if err != nil {
		t.Errorf("art2 error: %v\n", err)
	}

	r = NewReader(art2.Body)
	_, err = bufio.NewReader(r).WriteTo(b)

	if err != nil {
		if err == ErrBadCRC {
			t.Logf("expected crc: %d, was %d.\n", r.ExpectedCRC, r.CRC.Sum32())
		}
		t.Errorf("Article 2 error: %v\n", err)
	}

	fenc, err := ioutil.ReadFile("./test/joystick.jpg")

	if err != nil {
		t.Fatalf("Error loading file for comparison")
	}

	if !bytes.Equal(b.Bytes(), fenc) {
		log.Println("got", hex.Dump(b.Bytes()))
		log.Println("expected", fenc)
		t.Error("Multipart did not decode properly (decoded separately)")
	}
}
