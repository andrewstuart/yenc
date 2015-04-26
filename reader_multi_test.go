package yenc

import (
	"bufio"
	"bytes"
	"io/ioutil"
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
		t.Fatalf("Article 1 err: %v", err)
	}

	bufio.NewReader(NewReader(art.Body)).WriteTo(b)

	f2, err := os.Open("./test/00000021.ntx")

	if err != nil {
		t.Fatalf("f2 error")
	}

	art2, err := nntp.NewResponse(f2)

	if err != nil {
		t.Fatalf("art2 error: %v", err)
	}

	_, err = bufio.NewReader(NewReader(art2.Body)).WriteTo(b)

	if err != nil {
		t.Fatal(err)
	}

	fenc, err := ioutil.ReadFile("./test/joystick.jpg")

	if err != nil {
		t.Fatalf("message")
	}

	if !bytes.Equal(b.Bytes(), fenc) {
		t.Error("Multipart did not decode properly (decoded separately)")
	}
}
