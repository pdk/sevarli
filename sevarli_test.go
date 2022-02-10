package main

import "testing"

func TestUntabify(t *testing.T) {

	for i, test := range []struct{ input, output string }{
		{"a", "a"},
		{"a\tb", "a       b"},
		{"aaa\tb", "aaa     b"},
		{"\tx", "        x"},
		{"\t\tx", "                x"},
		{"aaaa\tbbbb\tcccc", "aaaa    bbbb    cccc"},
		{"aaaaaaa\tbbbbbbb\tccccccc\td", "aaaaaaa bbbbbbb ccccccc d"},
	} {

		result := untabify(test.input)
		if result != test.output {
			t.Errorf("%d: expected %#v, but got %#v", i, test.output, result)
		}

	}

}
