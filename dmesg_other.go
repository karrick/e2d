// +build !linux

package main

var dummy bool
var optDmesg *bool = &dummy

func getOffset() (int64, error) { return 0, nil }
