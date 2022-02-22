package main

import (
	"bufio"
	"io"
	"log"
	"os"

	"github.com/pdk/sevarli/comment"
	"github.com/pdk/sevarli/tabify"
	"github.com/pdk/sevarli/whitespace"
)

func main() {

	// fəˈtäɡrəfər
	// 写真家

	var r io.Reader

	r = tabify.NewTabsToSpacesReader(os.Stdin)
	r, err := comment.NewRemoveCommentsReader(r, "(^|[ \t])(#|//)")
	if err != nil {
		log.Fatalf("failed to create comment remover: %v", err)
	}
	r = whitespace.NewTrimRightWhiteSpaceReader(r)
	w := bufio.NewWriter(os.Stdout)
	for {

		buf := make([]byte, 3)
		b, err := r.Read(buf)
		if err != nil {
			return
		}
		w.WriteString(string(buf[0:b]))
		w.Flush()
	}
}
