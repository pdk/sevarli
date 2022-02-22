package comment

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
)

const defaultBufferSize = 4096

type RemoveCommentsReader struct {
	scanner               *bufio.Scanner
	bufferPos, bufferUsed int
	buffer                []byte
	pattern               *regexp.Regexp
	regexpError           error
}

func NewRemoveCommentsReader(input io.Reader, pattern string) (*RemoveCommentsReader, error) {

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	return &RemoveCommentsReader{
		scanner: bufio.NewScanner(input),
		buffer:  make([]byte, defaultBufferSize),
		pattern: regex,
	}, nil
}

func (ts *RemoveCommentsReader) Read(p []byte) (int, error) {
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

func (r *RemoveCommentsReader) readNextLine() (int, error) {

	ok := r.scanner.Scan()
	if !ok {
		return 0, io.EOF
	}
	if r.scanner.Err() != nil {
		return 0, r.scanner.Err()
	}

	line := r.scanner.Text()
	loc := r.pattern.FindStringIndex(line)
	if loc != nil {
		line = line[0:loc[0]]
	}

	if len(line)+1 > len(r.buffer) {
		return 0, fmt.Errorf("RemoveCommentsReader read buffer size %d is smaller than input line length %d", len(r.buffer), len(line))
	}

	copy(r.buffer, []byte(line))
	r.buffer[len(line)] = '\n'

	return len(line) + 1, nil
}
