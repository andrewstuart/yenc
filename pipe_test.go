package yenc

import (
	"io"
	"io/ioutil"
	"strings"
	"sync"
	"testing"
)

func TestFoo(t *testing.T) {
	r, w := io.Pipe()

	yr := NewReader(r)
	yw := NewWriter(w)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer w.Close()
		defer yw.Close()
		n, err := io.Copy(yw, strings.NewReader(a))

		if int(n) != len(a) {
			t.Errorf("Did not write all %d chars: %d written\n", len(a), n)
		}

		if err != nil {
			t.Fatalf("Could not copy. %v", err)
		}

		if err != nil {
			t.Fatalf("Error when closing yenc writer: %v", err)
		}

		wg.Done()
	}()

	bs, err := ioutil.ReadAll(yr)
	r.Close()

	if err != nil {
		t.Errorf("Error read all for pipe: %v", err)
	}

	if string(bs) != a {
		t.Errorf("Wrong string read: %s", string(bs))
	}

	wg.Wait()
}

const a = `
foo bar baz bang bang
`
