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

By default the output date-time strings are displayed in the local
time zone. This program can also display the date-time string in UTC
format by providing the --utc command line flag.

When provided command line arguments, this program converts each one
to a human readable date-time string. When no non-option command line
arguments are provided, this program reads from standard input, and
converts each line to the corresponding date-time string.

```Bash
$ e2d 1502400972
2017-08-10 17:36:12 -0400 EDT
```
