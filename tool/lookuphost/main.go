package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/jreisiger/stdintool"
)

func main() {
	stdintool.Run(&lookup{}, 1000)
}

type lookup struct {
	host string

	addrs []string
	err   error
}

func (*lookup) Parse(line string) stdintool.Task {
	return &lookup{host: line}
}

func (h *lookup) Process() {
	h.addrs, h.err = net.LookupHost(h.host)
}

func (h *lookup) Print() {
	if h.err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", filepath.Base(os.Args[0]), h.err)
		return
	}
	fmt.Printf("%s: %s\n", h.host, strings.Join(h.addrs, ", "))
}
