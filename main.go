package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [-s <seconds> | -m <milliseconds> | -n <nanoseconds>]\n", filepath.Base(os.Args[0]))
	os.Exit(2)
}

func main() {
	var seconds, milliseconds, nanoseconds float64

	flag.Float64Var(&seconds, "s", math.NaN(), "use seconds")
	flag.Float64Var(&milliseconds, "m", math.NaN(), "use milliseconds")
	flag.Float64Var(&nanoseconds, "n", math.NaN(), "use nanoseconds")
	flag.Parse()

	if flag.NArg() != 0 {
		usage()
	}

	if !math.IsNaN(seconds) {
		// no-op
	} else if !math.IsNaN(milliseconds) {
		seconds = milliseconds * float64(time.Millisecond) / float64(time.Second)
	} else if !math.IsNaN(nanoseconds) {
		seconds = nanoseconds * float64(time.Nanosecond) / float64(time.Second)
	} else {
		usage()
	}

	fmt.Println(time.Unix(int64(seconds), 0))
}
