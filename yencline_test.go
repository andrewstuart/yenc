package yenc

import "testing"

func TestHeader(t *testing.T) {
	a := &Header{}
	a.Put("size", "123")
	s := a.String()

	if s != "size=123" {
		t.Errorf("Wrong string: %s", s)
	}

	if a.Get("size") != "123" {
		t.Errorf("Wrong value returned: %s", a.Get("size"))
	}

	b := []byte("foo=3 bar=baz test=true\n")
	y, n := ReadYENCHeader(b)

	if len(*y) != 3 {
		t.Fatalf("Wrong number of headers\n")
	}

	if y.Get("foo") != "3" {
		t.Errorf("wrong value at foo\n")
	}

	if len(b) != n {
		t.Errorf("ReadYENCHeader did not modify original byte slice. Length: %d\n", n)
	}
}
