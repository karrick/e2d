package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/karrick/golf"
	"github.com/karrick/golinewrap"
	"github.com/karrick/gows"
)

var (
	optHelp         = golf.BoolP('h', "help", false, "When true, displays help then exits.")
	optMilliseconds = golf.BoolP('m', "milliseconds", false, "Use milliseconds rather than seconds.")
	optNanoseconds  = golf.BoolP('n', "nanoseconds", false, "Use nanoseconds rather than seconds.")
	optUTC          = golf.BoolP('u', "utc", false, "Display times in UTC rather than in local time zone.")
)

func help(err error) {
	lw := lineWrapping(os.Stderr)

	if err != nil {
		_, _ = fmt.Fprintf(lw, "ERROR: %s", err)
		_, _ = fmt.Fprintln(os.Stderr) // force additional newline after error message
	}

	_, _ = fmt.Fprintf(lw, "Simple CLI application to convert a epoch value to a date-time string.")

	golf.Usage()
	_, _ = fmt.Fprintf(os.Stderr, "\nUSAGE:\t%s [-m | -n] [-u] [epoch1 [epoch2 ...]]\n\n", filepath.Base(os.Args[0]))

	message := `With one or more command line arguments, displays the
        corresponding human readable time values. Without command line
        arguments, displays the corresponding human readable time value for
        each line of standard input.`
	_, _ = fmt.Fprintf(lw, message)
}

func main() {
	golf.Parse()

	if *optHelp {
		help(nil)
		os.Exit(0)
	}

	if *optMilliseconds && *optNanoseconds {
		help(errors.New("cannot provide both --milliseconds and --nanoseconds command line flags."))
		os.Exit(2)
	}

	divisor := 1.0
	if *optMilliseconds {
		divisor = float64(time.Second / time.Millisecond)
	} else if *optNanoseconds {
		divisor = float64(time.Second / time.Nanosecond)
	}

	if golf.NArg() > 0 {
		for _, arg := range golf.Args() {
			display(divisor, arg)
		}
		os.Exit(0)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		display(divisor, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}

func display(divisor float64, value string) {
	f64, err := strconv.ParseFloat(value, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: %s\n", err)
		return
	}

	f64 /= divisor // convert value to seconds

	sec := int64(f64)
	nsec := int64((f64 - float64(sec)) * float64(time.Second/time.Nanosecond))

	t := time.Unix(sec, nsec)
	if *optUTC {
		t = t.UTC()
	}

	fmt.Println(t.String())
}

func lineWrapping(w io.Writer) io.Writer {
	columns, _, err := gows.GetWinSize()
	if err != nil {
		return w
	}

	lw, err := golinewrap.New(w, columns, "")
	if err != nil {
		return w
	}

	return lw
}
