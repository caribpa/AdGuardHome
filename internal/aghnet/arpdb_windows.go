//go:build windows
// +build windows

package aghnet

import (
	"bufio"
	"io"
	"net"
	"strings"
	"sync"
)

func newARPDB() *cmdARPDB {
	return &cmdARPDB{
		runcmd: rcArpA,
		ns: &neighs{
			mu: &sync.RWMutex{},
			ns: make([]Neighbor, 0),
		},
		parse: parseArpA,
	}
}

// rcArpA runs "arp /a".
func rcArpA() (r io.Reader, err error) {
	return runCmd("arp", "/a")
}

// parseArpA parses the output of the "arp /a" command on Windows.  The expected
// input format (the first line is empty):
//
//
//   Interface: 192.168.56.16 --- 0x7
//     Internet Address      Physical Address      Type
//     192.168.56.1          0a-00-27-00-00-00     dynamic
//     192.168.56.255        ff-ff-ff-ff-ff-ff     static
//
func parseArpA(sc *bufio.Scanner, lenHint int) (ns []Neighbor) {
	ns = make([]Neighbor, 0, lenHint)
	for sc.Scan() {
		ln := sc.Text()
		if ln == "" {
			continue
		}

		fields := strings.Fields(ln)
		if len(fields) != 3 {
			continue
		}

		n := Neighbor{}

		if ip := net.ParseIP(fields[0]); ip == nil {
			continue
		} else {
			n.IP = ip
		}

		if mac, err := net.ParseMAC(fields[1]); err != nil {
			continue
		} else {
			n.MAC = mac
		}

		ns = append(ns, n)
	}

	return ns
}
