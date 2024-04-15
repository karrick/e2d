# e2d

Epoch to Date conversion.

## Description

Simple CLI application to convert a epoch value to a date-time string.

## Usage

By default input values are assumed to be provided in the number of
seconds since the epoch. This program can also convert the number of
milliseconds since the epoch, or the number of nanoseconds since the
epoch by providing either the --milliseconds or the --nanoseconds
command line flag.

By default this program displays the date-time strings in your
computer's local time zone, which may or may not be UTC. To force this
program to output the date-time strings in UTC, use the --utc command
line flag. If your computer is configured to use UTC by default, there
is no way for this program to know what your local time zone is, so
output will always be in UTC.

When provided command line arguments, this program converts each
to a human readable date-time string. When no non-option command line
arguments are provided, this program reads from standard input, and
converts each line to the corresponding date-time string.

```Bash
$ e2d 1502400972
2017-08-10 17:36:12 -0400 EDT
```

### Translate the value of a particular text field.

```Bash
$ e2d --field 3 space-delimited.log
```

### Translate the value of a particular text field with a non-space delimiter.

```Bash
$ e2d --microseconds --delimiter , --field 3 comma-delimited.txt
```

### Translate the value of a JSON property.

```Bash
$ e2d --nanoseconds --property epoch newline-terminated-json.log
```

## Installation

If you don't have the Go programming language installed, then you'll
need to install a copy from
[https://golang.org/dl](https://golang.org/dl).

Once you have Go installed:

    $ go install github.com/karrick/e2d@latest
