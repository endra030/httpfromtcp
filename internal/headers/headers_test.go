package headers

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test: Valid single header

func TestHeaders(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Valid single header with extra whitespace
	headers = NewHeaders()
	data = []byte("Accept-Encoding:    gzip    \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "gzip", headers["Accept-Encoding"])
	expectedParsed := bytes.Index(data, []byte("\r\n")) + 2
	assert.Equal(t, expectedParsed, n)
	assert.False(t, done)

	// Test: Valid 2 headers with existing headers
	headers = NewHeaders()
	data = []byte("Content-Type: application/json\r\nContent-Length: 123\r\n\r\n")
	total := 0

	// First header
	n, done, err = headers.Parse(data[total:])
	require.NoError(t, err)
	assert.Equal(t, "application/json", headers["Content-Type"])
	total += n
	assert.False(t, done)

	// Second header
	n, done, err = headers.Parse(data[total:])
	require.NoError(t, err)
	assert.Equal(t, "123", headers["Content-Length"])
	total += n
	assert.False(t, done)

	// Final empty line (end of headers)
	n, done, err = headers.Parse(data[total:])
	require.NoError(t, err)
	total += n
	assert.True(t, done)

	// Total bytes parsed should equal the length of full data
	assert.Equal(t, len(data), total)

	// Test: Valid done
	headers = NewHeaders()
	data = []byte("Host: example.com\r\n\r\n")
	total = 0

	// First header
	n, done, err = headers.Parse(data[total:])
	require.NoError(t, err)
	assert.Equal(t, "example.com", headers["Host"])
	total += n
	assert.False(t, done)

	// Final empty line (end of headers)
	n, done, err = headers.Parse(data[total:])
	require.NoError(t, err)
	total += n
	assert.True(t, done)

	// Validate full consumption
	assert.Equal(t, len(data), total)

}
