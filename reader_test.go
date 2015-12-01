package yenc

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestDecoder(t *testing.T) {
	encoded, err := os.Open("./test/00000005.ntx")

	if err != nil {
		t.Fatalf("could not begin test: ntx file not readable")
	}

	unencoded, err := ioutil.ReadFile("./test/testfile.txt")

	if err != nil {
		t.Fatalf("could not read unencoded test file")
	}

	decBytes, err := ioutil.ReadAll(NewReader(encoded))

	if err != nil {
		t.Fatalf("Error reading decoded bytes: %v", err)
	}

	if !bytes.Equal(decBytes, unencoded) {
		fmt.Println(hex.Dump(decBytes))

		diff := getDiff(unencoded, decBytes)
		t.Errorf("Decoded bytes did not equal unencoded bytes. Diff was %d long; dec was %d long", len(diff), len(decBytes))
	}
}

func BenchmarkDecoder(b *testing.B) {
	encoded, err := os.Open("./test/00000005.ntx")

	if err != nil {
		b.Fatalf("could not begin test: ntx file not readable")
	}

	bs, err := ioutil.ReadAll(encoded)
	if err != nil {
		log.Fatal(err)
	}

	buf := bytes.NewBuffer(bs)

	empty := make([]byte, 2<<20)

	for i := 0; i < b.N; i++ {
		NewReader(buf).Read(empty)
	}
}

func getDiff(b1, b2 []byte) []byte {
	b1c := make([]byte, len(b1))
	copy(b1c, b1)

	b2c := make([]byte, len(b2))
	copy(b2c, b2)

	for len(b1c) > 0 && len(b2c) > 0 && b1c[0] == b2c[0] {
		b1c = b1c[1:]
		b2c = b2c[1:]
	}

	if len(b1c) > 0 {
		return b1c
	}

	return b2c
}
