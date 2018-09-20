package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/karrick/golf"
	"github.com/karrick/golinewrap"
	"github.com/karrick/gows"
)

var (
	optDelimiter    = golf.StringP('d', "delimiter", " ", "Use STRING as field delimiter rather than whitespace.")
	optField        = golf.UintP('f', "field", 0, "When not 0, specifies the field number to convert to a time.")
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
	_, _ = fmt.Fprintf(os.Stderr, "\nUSAGE:\t%s [-m | -n] [-u] [--uptime] [epoch1 [epoch2 ...]]\n\n", filepath.Base(os.Args[0]))

	message := `With one or more command line arguments, displays the
        corresponding human readable time values. Without command line
        arguments, program runs as a filter, attempting to convert either the
        entire line, or a specified field, from a numerical epoch to a human
        readable date value. When running as a filter, if the line, or specified
        field, fails to parse as a number, this program emits that line
        unchanged and suppresses all warning messages.`

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

	offset, err := getOffset()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}

	if golf.NArg() > 0 {
		for _, arg := range golf.Args() {
			s, err := display(divisor, offset, arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "WARNING: %s\n", err)
				continue
			}
			fmt.Println(s)
		}
		os.Exit(0)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if !*optDmesg && *optField == 0 {
			s, err := display(divisor, offset, line)
			if err != nil {
				fmt.Println(line)
				continue
			}
			fmt.Println(s)
			continue
		}

		if len(line) == 0 {
			fmt.Println()
			continue
		}

		if *optDmesg {
			// first character is `[`
			if line[0] != '[' {
				fmt.Println(line)
				continue
			}

			// grab all until `]`
			end := strings.Index(line, "]")
			if end < 0 {
				fmt.Println(line)
				continue
			}

			start := strings.LastIndex(line[:end], " ") + 1
			if start == 0 {
				// When cannot find first space, then there is none, and should
				// start at first index, which is the character immediately
				// following the [.
				start = 1
			}

			s, err := display(divisor, offset, line[start:end])
			if err != nil {
				// could not convert field, so emit unchanged line
				fmt.Println(line)
				continue
			}

			fmt.Printf("[%s%s\n", s, line[end:])
			continue
		}

		var fields []string
		if *optDelimiter == " " {
			fields = strings.Fields(line)
		} else {
			fields = strings.Split(line, *optDelimiter)
		}

		if uint(len(fields)) < *optField {
			// line does not have enough fields, so emit unchanged line
			fmt.Println(line)
			continue
		}

		s, err := display(divisor, offset, fields[*optField-1])
		if err != nil {
			// could not convert field, so emit unchanged line
			fmt.Println(line)
			continue
		}

		fields[*optField-1] = s
		fmt.Println(strings.Join(fields, *optDelimiter))
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}

func display(divisor float64, offset int64, value string) (string, error) {
	f64, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return "", err
	}

	f64 /= divisor // convert value to seconds

	sec := int64(f64)
	nsec := int64((f64 - float64(sec)) * float64(time.Second/time.Nanosecond))

	t := time.Unix(sec+offset, nsec)
	if *optUTC {
		t = t.UTC()
	}

	return t.String(), nil
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
