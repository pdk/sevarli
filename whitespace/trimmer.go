package whitespace

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"
)

// defaultBufferSize is the longest line we can handle
const defaultBufferSize = 4096

type TrimRightWhiteSpaceReader struct {
	scanner               *bufio.Scanner
	bufferPos, bufferUsed int
	buffer                []byte
}

// NewTrimRightWhiteSpaceReader returns a new io.Reader which will trim white
// space from the right side of input lines.
func NewTrimRightWhiteSpaceReader(input io.Reader) *TrimRightWhiteSpaceReader {

	return &TrimRightWhiteSpaceReader{
		scanner: bufio.NewScanner(input),
		buffer:  make([]byte, defaultBufferSize),
	}
}

// Read processes input, removing trailing whitespace.
func (ts *TrimRightWhiteSpaceReader) Read(p []byte) (int, error) {
	var err error

	if ts.bufferPos >= ts.bufferUsed {
		ts.bufferPos = 0
		ts.bufferUsed, err = ts.readNextLine()
		if err != nil {
			return 0, err
		}
	}

	avail := ts.bufferUsed - ts.bufferPos
	c := len(p)
	if avail < c {
		c = avail
	}

	copy(p, ts.buffer[ts.bufferPos:ts.bufferPos+c])
	ts.bufferPos += c

	return c, nil
}

func (r *TrimRightWhiteSpaceReader) readNextLine() (int, error) {

	ok := r.scanner.Scan()
	if !ok {
		return 0, io.EOF
	}
	if r.scanner.Err() != nil {
		return 0, r.scanner.Err()
	}

	line := strings.TrimRightFunc(r.scanner.Text(), unicode.IsSpace)

	if len(line)+1 > len(r.buffer) {
		return 0, fmt.Errorf("TrimRightWhiteSpaceReader read buffer size %d is smaller than input line length %d", len(r.buffer), len(line))
	}

	copy(r.buffer, []byte(line))
	r.buffer[len(line)] = '\n'

	return len(line) + 1, nil
}
