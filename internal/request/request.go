package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	dataBytes, err := io.ReadAll(reader)
	req := Request{}
	if err != nil {
		return &req, err

	}
	dataString := string(dataBytes)
	dataSlice := strings.Split(dataString, "\r\n")
	reqLineString := dataSlice[0]
	reqLineSlice := strings.Split(reqLineString, " ")
	fmt.Printf("dataSlice:%v\n", dataSlice)
	if len(reqLineSlice) != 3 {
		return &req, errors.New("invalid request line")
	}

	for _, r := range reqLineSlice[0] {
		if !unicode.IsLetter(r) {
			return &req, errors.New("invalid http method")

		}
		if !unicode.IsUpper(r) {
			return &req, errors.New("invalid http method")
		}
	}

	httpVersion := strings.Split(reqLineSlice[2], "/")[1]
	if httpVersion != "1.1" {
		return &req, errors.New("unsupported http version")
	}
	req.RequestLine.HttpVersion = httpVersion
	req.RequestLine.Method = reqLineSlice[0]
	req.RequestLine.RequestTarget = reqLineSlice[1]
	return &req, nil

}
