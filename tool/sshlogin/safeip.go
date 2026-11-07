package main

import (
	"errors"
	"net"
	"strings"
)

type safeIPaddresses []string

func (ip *safeIPaddresses) String() string { return strings.Join(*ip, ",") }
func (ip *safeIPaddresses) Set(v string) error {
	addr := net.ParseIP(v)
	if addr == nil {
		return errors.New("not an IP address")
	}
	*ip = append(*ip, v)
	return nil
}
