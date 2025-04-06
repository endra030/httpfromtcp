package headers

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return 0, false, nil
	}
	line := string(data[:idx])
	fmt.Printf("value of line= %s\n", line)

	if line == "" {
		return idx + 2, true, nil // empty line → end of headers
	}
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return 0, false, errors.New("malformed header line")
	}
	rawKey := parts[0]
	rawValue := parts[1]

	if strings.TrimSpace(rawKey) != rawKey {
		return 0, false, errors.New("invalid spacing in header key")
	}
	key := rawKey
	value := strings.TrimSpace(rawValue)
	if key == "" || value == "" {
		return 0, false, errors.New("empty header key or value")
	}
	h[key] = value
	return idx + 2, false, nil

}
