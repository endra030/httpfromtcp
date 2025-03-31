package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
	ParserState ParserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type ParserState int

const (
	initialized ParserState = iota
	done
)
const bufferSize = 8

// func RequestFromReader(reader io.Reader) (*Request, error) {
// 	buf := make([]byte, bufferSize, bufferSize)
// 	readToIndex := 0
// 	req := Request{
// 		ParserState: initialized,
// 	}
// 	for req.ParserState != done {
// 		if readToIndex >= len(buf) {
// 			newBuf := make([]byte, len(buf)*2)
// 			copy(newBuf, buf)
// 			buf = newBuf
// 		}
// 		numBytesRead, err := reader.Read(buf[readToIndex:])
// 		if err != nil {
// 			if errors.Is(io.EOF, err) {
// 				req.ParserState = done
// 				break
// 			}
// 			return nil, err
// 		}
// 		readToIndex += numBytesRead

// 		numBytesParsed, err := req.parse(buf[:readToIndex])
// 		if err != nil {
// 			return nil, err
// 		}

// 		copy(buf, buf[numBytesParsed:])
// 		readToIndex -= numBytesParsed
// 	}

// 	return &req, nil

// }

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0
	req := &Request{
		ParserState: initialized,
	}
	for req.ParserState != done {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		numBytesRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(io.EOF, err) {
				req.ParserState = done
				break
			}
			return nil, err
		}
		readToIndex += numBytesRead

		numBytesParsed, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[numBytesParsed:])
		readToIndex -= numBytesParsed
	}
	return req, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	fmt.Printf("parseRequestLineData:%s\n", string(data))
	reqLine := &RequestLine{}
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return reqLine, 0, nil
	}
	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return reqLine, 0, err
	}
	reqLine.HttpVersion = requestLine.HttpVersion
	reqLine.Method = requestLine.Method
	reqLine.RequestTarget = requestLine.RequestTarget
	return reqLine, idx + 2, nil

}

func requestLineFromString(dataString string) (*RequestLine, error) {
	req := RequestLine{}
	dataSlice := strings.Split(dataString, "\r\n")
	reqLineString := dataSlice[0]
	reqLineSlice := strings.Split(reqLineString, " ")
	fmt.Printf("dataSlice:%v\n", dataSlice)
	if len(reqLineSlice) != 3 {
		return nil, errors.New("invalid request line")
	}

	for _, r := range reqLineSlice[0] {
		if !unicode.IsLetter(r) {
			return nil, errors.New("invalid http method")

		}
		if !unicode.IsUpper(r) {
			return nil, errors.New("invalid http method")
		}
	}

	httpVersion := strings.Split(reqLineSlice[2], "/")[1]
	if httpVersion != "1.1" {
		return nil, errors.New("unsupported http version")
	}
	req.HttpVersion = httpVersion
	req.Method = reqLineSlice[0]
	req.RequestTarget = reqLineSlice[1]
	return &req, nil

}

func (r *Request) parse(data []byte) (int, error) {

	if r.ParserState == initialized {
		fmt.Printf("parseData:%s\n", string(data))
		reqL, n, err := parseRequestLine(data)
		if err != nil {

			return n, err
		}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = *reqL
		r.ParserState = done
		return n, nil

	} else if r.ParserState == done {
		return 0, errors.New("error: trying to read data in a done state")
	} else {
		return 0, errors.New("unknown state")

	}

}

// func (r *Request) parse(data []byte) (int, error) {
// 	switch r.ParserState {
// 	case initialized:
// 		fmt.Printf("parseData:%s\n", string(data))
// 		requestLine, n, err := parseRequestLine(data)
// 		if err != nil {
// 			// something actually went wrong
// 			return 0, err
// 		}
// 		if n == 0 {
// 			// just need more data
// 			return 0, nil
// 		}
// 		r.RequestLine = *requestLine
// 		r.ParserState = done
// 		return n, nil
// 	case done:
// 		return 0, fmt.Errorf("error: trying to read data in a done state")
// 	default:
// 		return 0, fmt.Errorf("unknown state")
// 	}
// }
