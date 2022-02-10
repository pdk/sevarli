package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

func main() {

	var args CommandArgs
	flag.StringVar(&args.data, "data", "", "path to data file (otherwise read stdin)")
	flag.BoolVar(&args.export, "export", false, "export the vars")
	flag.StringVar(&args.pattern, "pattern", "", "pattern to search for (*required)")
	flag.Var(&args.hide, "hide", "hide column(s) when listing")
	flag.StringVar(&args.prefix, "prefix", "", "prefix variable name with given value")
	flag.StringVar(&args.suffix, "suffix", "", "suffix variable name with given value")
	flag.BoolVar(&args.caps, "caps", true, "convert names to caps")

	flag.Parse()

	if args.pattern == "" {
		fmt.Fprintf(os.Stderr, "usage of sevarli:\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if err := run(args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(args CommandArgs) error {

	data, err := readLines(args.data)
	if err != nil {
		return err
	}

	data = stripBlankLines(stripComments(untabifies(data)))
	cols := identifyColumns(data)

	matched, err := matching(data[1:], args.pattern)
	if err != nil {
		return err
	}

	if len(matched) == 0 {
		printOptions(data[0], data[1:])
		return fmt.Errorf("did not match any line")
	}

	if len(matched) != 1 {
		printOptions(data[0], matched)
		return fmt.Errorf("matched more than 1 line")
	}

	for _, l := range matched {
		setCols(cols, l, args.export, args.prefix, args.suffix, args.caps)
	}

	return nil
}

func printOptions(header string, lines []string) {
	fmt.Fprintln(os.Stderr, header)
	for _, l := range lines {
		fmt.Fprintln(os.Stderr, l)
	}
}

func setCols(cols columns, data string, export bool, prefix, suffix string, caps bool) {
	exp := ""
	if export {
		exp = "export "
	}

	datar := []rune(data)

	for _, col := range cols {
		name := fixName(col.name, caps)
		r := col.last + 1
		if r > len(datar) {
			r = len(datar)
		}
		value := strings.TrimSpace(string(datar[col.first:r]))
		fmt.Printf("%s%s%s%s=\"%s\"\n", exp, prefix, name, suffix, value)
	}
}

func fixName(name string, caps bool) string {

	if caps {
		name = strings.ToUpper(name)
	}

	return strings.ReplaceAll(name, " ", "_")
}

func matching(lines []string, pattern string) ([]string, error) {

	matcher, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("pattern %#v is invalid: %w", pattern, err)
	}

	newLines := make([]string, 0, len(lines))

	for _, l := range lines {
		if matcher.MatchString(l) {
			newLines = append(newLines, l)
		}
	}

	return newLines, nil
}

func stripBlankLines(lines []string) []string {
	newLines := make([]string, 0, len(lines))

	for _, l := range lines {
		if l != "" {
			newLines = append(newLines, l)
		}
	}

	return newLines
}

func stripComments(lines []string) []string {

	for i, line := range lines {
		line = removeComment(line, "#")
		line = removeComment(line, "//")
		lines[i] = line
	}

	return lines
}

func removeComment(input, marker string) string {
	if strings.HasPrefix(input, marker) {
		return ""
	}

	// only recognize marker if having leading whitespace
	i := strings.Index(input, " "+marker)
	if i < 0 {
		return input
	}

	return strings.TrimRight(input[:i], " ")
}

// untabifies untabifies a slice of strings
func untabifies(lines []string) []string {
	for i, line := range lines {
		lines[i] = untabify(line)
	}

	return lines
}

// untabify replaces any tab characters with the right number of spaces
func untabify(s string) string {

	i := strings.IndexRune(s, '\t')
	if i < 0 {
		return strings.TrimRight(s, " ")
	}

	newS := s[:i] + strings.Repeat(" ", 8-(i%8)) + s[i+1:]

	return untabify(newS)
}

func readLines(path string) ([]string, error) {

	if path == "" {
		return readInput(os.Stdin)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %w", err)
	}
	defer file.Close()

	return readInput(file)
}

func readInput(file io.Reader) ([]string, error) {
	var lines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return lines, fmt.Errorf("failed to read data: %w", err)
	}

	return lines, nil
}
