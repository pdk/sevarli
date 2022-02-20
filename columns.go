package main

import (
	"log"
	"unicode/utf8"
)

type column struct {
	name  string
	first int // first index of column (by rune)
	last  int // last index+1 of column (by rune)
}

func (c column) width() int {
	return c.last - c.first
}

type columns []column

func identifyColumns(lines []string) columns {

	nonSpace := make([]bool, maxLen(lines))

	// across all lines, identify which columns have non-space characters.
	for _, l := range lines {
		i := 0
		for _, c := range l {
			if c != ' ' {
				nonSpace[i] = true
			}
			i++
		}
	}

	// mark any is-space column as not-is-space if only a single space.
	// X.X => XXX
	// X..X = X..X
	for i, b := range nonSpace {
		if !b && i > 0 && i < len(nonSpace)-1 && nonSpace[i-1] && nonSpace[i+1] {
			nonSpace[i] = true
		}
	}

	// identify first, last position of each "column".
	// "last" position is index of first space after column.
	firsts, lasts := []int{}, []int{}
	wasNonSpace := false
	for i := 0; i < len(nonSpace); i++ {
		switch {
		case !wasNonSpace && nonSpace[i]:
			firsts = append(firsts, i)
		case wasNonSpace && !nonSpace[i]:
			lasts = append(lasts, i)
		}
		wasNonSpace = nonSpace[i]
	}
	if len(firsts) > len(lasts) {
		lasts = append(lasts, len(nonSpace))
	}
	if len(firsts) != len(lasts) {
		// bug should never happen
		log.Fatalf("failed to identify columns for %#v", nonSpace)
	}

	// pull out names for columns from the first data line
	headerLine := []rune(lines[0])
	cols := columns{}
	for i := 0; i < len(firsts); i++ {
		cols = append(cols, column{
			first: firsts[i],
			last:  lasts[i],
			name:  getColVal(headerLine, firsts[i], lasts[i]),
		})
	}

	return cols
}

func maxLen(lines []string) int {

	m := 0
	for _, l := range lines {
		len := utf8.RuneCountInString(l)
		if len > m {
			m = len
		}
	}

	return m
}
