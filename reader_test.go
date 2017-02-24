package yenc

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecoder(t *testing.T) {
	asrt := assert.New(t)
	encoded, err := ioutil.ReadFile("./test/00000005.ntx")
	asrt.NoError(err)

	unencoded, err := ioutil.ReadFile("./test/testfile.txt")
	asrt.NoError(err)

	decBytes, err := ioutil.ReadAll(NewReader(bytes.NewReader(encoded)))
	asrt.NoError(err)

	if !asrt.Equal(decBytes, unencoded) {
		fmt.Println("predecoded")
		fmt.Println(hex.Dump(encoded))
		fmt.Println("decoded")
		fmt.Println(hex.Dump(decBytes))
		fmt.Println("expected")
		fmt.Println(hex.Dump(unencoded))

		diff := getDiff(unencoded, decBytes)
		t.Errorf("Decoded bytes did not equal unencoded bytes. Diff was %d long; dec was %d long", len(diff), len(decBytes))
	}
}

func TestReader(t *testing.T) {
	asrt := assert.New(t)
	fs, _ := filepath.Glob("./test/examples/*raw")

	for _, e := range fs {
		enc, err := os.Open(e)
		asrt.NoError(err)

		r := NewReader(enc)

		b := &bytes.Buffer{}
		_, err = io.Copy(b, r)
		asrt.NoError(err)
	}
}

func BenchmarkReader(b *testing.B) {
	encoded, err := ioutil.ReadFile("./test/00000005.ntx")

	if err != nil {
		b.Fatalf("could not begin test: ntx file not readable")
	}

	b.SetBytes(int64(len(encoded)))

	b.RunParallel(func(pb *testing.PB) {
		empty := make([]byte, 1024)

		for pb.Next() {
			r := NewReader(bytes.NewReader(encoded))
			n, err := r.Read(empty)

			if err != nil && err != io.EOF {
				b.Fatal(err)
			}
			if n == 0 {
				log.Println(empty)
				b.Fatal("Didn't read anything")
			}
		}
	})
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
