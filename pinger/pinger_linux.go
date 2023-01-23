package pinger

import (
	"os/exec"
	"strings"
)

// Pinger - Currently used below flags that are used among common Linux distros
// -c <count>         stop after <count> replies
// -w <deadline>      reply wait <deadline> in seconds
func (t Target) Pinger() (status bool) {
	out, _ := exec.Command("ping", t.Address, "-c", t.Count, "-w", t.Deadline).Output()
	if strings.Contains(string(out), "Destination Host Unreachable") {
		return false
	}
	return true
}
