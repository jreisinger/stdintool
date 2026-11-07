/*
ssh example.net "sudo journalctl --unit ssh --since '1 month ago'" | \
go run ./tool/sshlogin -s $(curl -s ifconfig.me) 2> /tmp/e
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/jreisiger/stdintool"
)

var (
	safeIPs safeIPaddresses
	tool    = filepath.Base(os.Args[0])
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", tool)
		flag.PrintDefaults()
	}
	flag.Var(&safeIPs, "s", "safe IP address to login (repeatable)")
	flag.Parse()

	stdintool.Run(&login{}, 1000)
}

type login struct {
	log string

	timestamp time.Time
	accepted  bool
	username  string
	method    string
	fromIP    string
	ipSafe    bool
	err       error
}

func (login) Parse(line string) stdintool.Task {
	return &login{log: line}
}

func (l *login) Process() {
	// 0   1  2		   3		   4			 5		  6			7   8    9    10
	// Nov 06 16:01:00 example.com sshd[552445]: Accepted publickey for root from 1.1.1.1 ...
	fields := strings.Fields(l.log)

	if len(fields) < 11 {
		l.err = fmt.Errorf("non-connection log: %q", l.log)
		return
	}

	if strings.ToLower(fields[5]) != "accepted" {
		// l.err = fmt.Errorf("non-accepted log: %q", l.log)
		return
	}
	l.accepted = true

	ts := strings.Join(fields[:3], " ")
	t, err := time.Parse("Jan 02 15:04:05", ts)
	if err != nil {
		l.err = err
		return
	}
	l.timestamp = t

	l.method = fields[6]
	l.username = fields[8]
	l.fromIP = fields[10]

	l.ipSafe = slices.Contains(safeIPs, l.fromIP)
}

func (l *login) Print() {
	if l.err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", tool, l.err)
		return
	}
	if l.accepted && !l.ipSafe {
		fmt.Printf("%s %s logged in via %s from %s\n", l.timestamp, l.username, l.method, l.fromIP)
	}
}
