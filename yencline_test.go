package yenc

import (
	"bufio"
	"strings"
	"testing"
)

func TestHeader(t *testing.T) {
	a := &YENCHeader{}
	a.Add("size", "123")
	s := a.String()

	if s != "size=123" {
		t.Errorf("Wrong string: %s", s)
	}

	if a.Get("size") != "123" {
		t.Errorf("Wrong value returned: %s", a.Get("size"))
	}

	b := bufio.NewReader(strings.NewReader("foo=3 bar=baz test=true\n"))
	y, err := ReadYENCHeader(b)

	if err != nil {
		t.Errorf("error reading header: %v", err)
	}

	if len(*y) != 3 {
		t.Fatalf("Wrong number of headers")
	}

	if y.Get("foo") != "3" {
		t.Errorf("wrong value at foo")
	}

}
