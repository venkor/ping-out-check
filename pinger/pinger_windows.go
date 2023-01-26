package pinger

import (
	"os/exec"
	"strconv"
	"strings"
)

// Pinger - Currently used below flags that are used in Windows 10
// -n count       Number of echo requests to send.
// -w timeout     Timeout in milliseconds to wait for each reply.
func (t Target) Pinger() (status bool) {
	out, _ := exec.Command("ping", t.Address, "-n", strconv.Itoa(t.Count), "-w", strconv.Itoa(t.Deadline)).Output()
	if strings.Contains(string(out), "Destination Host Unreachable") {
		return false
	}
	return true
}
