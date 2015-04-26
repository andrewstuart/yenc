package yenc

import (
	"bufio"
	"fmt"
	"strings"
)

type YENCHeader map[string]string

func (y *YENCHeader) String() string {
	s := make([]string, 0)

	for k, v := range *y {
		s = append(s, fmt.Sprintf("%s=%s", k, v))
	}

	return strings.Join(s, " ")
}

func (y *YENCHeader) Add(k, v string) {
	(*y)[k] = v
}

func (y *YENCHeader) Get(k string) string {
	return (*y)[k]
}

func ReadYENCHeader(br *bufio.Reader) (*YENCHeader, error) {
	s, err := br.ReadString('\n')

	if err != nil {
		return nil, err
	}

	s = strings.TrimSpace(s)

	y := &YENCHeader{}
	if err != nil {
		return y, err
	}

	ss := strings.Split(s, " ")
	for _, kvString := range ss {
		kvPair := strings.Split(kvString, "=")

		if len(kvPair) == 2 {
			y.Add(kvPair[0], kvPair[1])
		}
	}

	return y, nil
}
