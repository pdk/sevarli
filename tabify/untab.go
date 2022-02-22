package tabify

import (
	"bufio"
	"fmt"
	"io"
	"unicode/utf8"
)

// defaultBufferSize is the longest line we can handle
const defaultBufferSize = 4096

// TabsToSpacesReader is a io.Reader that converts tabs to spaces. Note: Use
// NewTabsToSpacesReader() to get an actually working instance.
type TabsToSpacesReader struct {
	scanner               *bufio.Scanner
	bufferPos, bufferUsed int
	buffer                []byte
}

// NewTabsToSpacesReader returns a Reader than translates any tabs ("\t") into
// the appropriate number of spaces. "Tabs" are every 8 spaces.
func NewTabsToSpacesReader(input io.Reader) *TabsToSpacesReader {

	return &TabsToSpacesReader{
		scanner: bufio.NewScanner(input),
		buffer:  make([]byte, defaultBufferSize),
	}
}

// Read reads the next line and converts tabs to spaces.
func (ts *TabsToSpacesReader) Read(p []byte) (int, error) {
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

func (ts *TabsToSpacesReader) readNextLine() (int, error) {

	ok := ts.scanner.Scan()
	if !ok {
		return 0, io.EOF
	}
	if ts.scanner.Err() != nil {
		return 0, ts.scanner.Err()
	}

	rbytes := make([]byte, 6)
	wp := 0
	line := ts.scanner.Text()

	cc := 0
	for _, r := range line {

		if wp >= len(ts.buffer) {
			return 0, fmt.Errorf("TabsToSpacesReader read buffer size %d is too small for input line length %d",
				len(ts.buffer), len(line))
		}

		if r == '\t' {
			ts.buffer[wp] = ' '
			wp++
			cc++
			for cc%8 != 0 {
				ts.buffer[wp] = ' '
				wp++
				cc++
			}
			continue
		}

		bc := utf8.EncodeRune(rbytes, r)
		for i := 0; i < bc; i++ {
			ts.buffer[wp] = rbytes[i]
			wp++
		}

		cc++
	}
	ts.buffer[wp] = '\n'
	wp++

	return wp, nil
}
