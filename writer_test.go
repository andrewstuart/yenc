package yenc

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestWrite(t *testing.T) {
	testWrite := &bytes.Buffer{}

	test, err := ioutil.ReadFile("./test/testfile.txt")

	if err != nil {
		t.Fatalf("Error reading file for test setup")
	}

	f, err := os.Create("./test/test.out")
	if err != nil {
		t.Fatalf("Error creating out file")
	}
	defer f.Close()

	mw := io.MultiWriter(f, testWrite)

	w := NewWriter(mw)
	w.Name = "testfile.txt"

	_, err = io.Copy(w, bytes.NewReader(test))

	if err != nil {
		t.Errorf("Error copying to writer: %v", err)
	}

	err = w.Close()

	if err != nil {
		t.Fatalf("Error closing: %v", err)
	}

	expected, err := ioutil.ReadFile("./test/00000005.nh.ntx")

	if err != nil {
		t.Fatalf("Error reading encoded file for test comparison: %v", err)
	}

	ourBytes := stripLines(testWrite.Bytes())
	expected = stripLines(expected)

	fmt.Printf("ourBytes = %s\n", ourBytes)

	if !bytes.Equal(ourBytes, expected) {
		t.Errorf("Wrong encoding produced.")
	}
}

func stripLines(b []byte) []byte {
	bs := bytes.Split(b, []byte("\n"))
	bs = bs[1:]
	bs = bs[:len(bs)-2]
	return bytes.Join(bs, []byte("\n"))
}

func TestIsWriter(t *testing.T) {
	takesWriter(NewWriter(nil))
}

func takesWriter(w io.Writer) {
}
