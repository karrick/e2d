// +build linux

package main

import (
	"syscall"
	"time"

	"github.com/karrick/golf"
)

var optDmesg = golf.Bool("dmesg", false, "Process input stream assuming it was produced by `dmesg`. (Only available on Linux.)")

func getOffset() (int64, error) {
	if !*optDmesg {
		return 0, nil
	}

	var info syscall.Sysinfo_t
	if err := syscall.Sysinfo(&info); err != nil {
		return 0, err
	}

	return time.Now().Unix() - info.Uptime, nil
}
