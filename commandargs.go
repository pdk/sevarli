package main

import (
	"strings"
)

type Strings []string

type CommandArgs struct {
	export  bool
	hide    Strings
	pattern string
	data    string
	prefix  string
	suffix  string
	caps    bool
}

func (s *Strings) Set(value string) error {

	for _, next := range strings.Split(value, ",") {
		*s = append(*s, strings.TrimSpace(next))
	}

	return nil
}

func (s *Strings) String() string {
	return strings.Join(*s, ",")
}
