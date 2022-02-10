package main

import (
	"strings"
	"unicode/utf8"
)

type column struct {
	name        string
	first, last int
}

type columns []column

func identifyColumns(lines []string) columns {

	nonSpace := make([]bool, maxLen(lines))

	for _, l := range lines {
		i := 0
		for _, c := range l {
			if c != ' ' && !nonSpace[i] {
				nonSpace[i] = true
			}
			i++
		}
	}

	for i, b := range nonSpace {
		if !b && i > 0 && nonSpace[i-1] && i < len(nonSpace) && nonSpace[i+1] {
			nonSpace[i] = true
		}
	}

	headerLine := []rune(lines[0])
	cols := columns{}
	isSpace := true
	first := 0
	for p := 0; p < len(nonSpace); p++ {
		if isSpace && nonSpace[p] {
			first = p
			isSpace = false
			continue
		}
		if !isSpace && (!nonSpace[p] || p == len(nonSpace)-1) {
			last := p - 1
			if p == len(nonSpace)-1 {
				last = p
			}
			col := column{
				name:  strings.TrimSpace(string(headerLine[first : last+1])),
				first: first,
				last:  last,
			}
			cols = append(cols, col)
			isSpace = true
		}
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
