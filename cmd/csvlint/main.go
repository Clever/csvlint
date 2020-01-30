package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"unicode/utf8"

	"github.com/Clever/csvlint"
)

func printHelpAndExit(code int) {
	flag.PrintDefaults()
	os.Exit(code)
}

func main() {
	delimiter := flag.String("delimiter", ",", "field delimiter in the file, for instance '\\t' or '|'")
	lazyquotes := flag.Bool("lazyquotes", false, "try to parse improperly escaped quotes")
	verbose := flag.Bool("verbose", false, "verbose output")
	help := flag.Bool("help", false, "print help and exit")
	flag.Parse()

	if *help {
		printHelpAndExit(0)
	}

	if flag.NFlag() > 0 {
		fmt.Fprintln(os.Stderr, "Warning: not using defaults, may not validate CSV to RFC 4180")
	}

	convertedDelimiter, err := strconv.Unquote(`'` + *delimiter + `'`)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error unquoting delimiter '%s', note that only one-character delimiters are supported\n", *delimiter)
		printHelpAndExit(1)
	}
	// don't need to check size since Unquote returns one-character string
	comma, _ := utf8.DecodeRuneInString(convertedDelimiter)

	if flag.NArg() > 1 {
		fmt.Fprint(os.Stderr, "csvlint accepts stdin or a single filepath as an argument\n")
		printHelpAndExit(1)
	}

	var f *os.File

	if flag.NArg() == 0 {
		f = os.Stdin
	} else {
		f, err = os.Open(flag.Arg(0))
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "file '%s' does not exist\n", flag.Arg(0))
				os.Exit(1)
			} else {
				panic(err)
			}
		}
		defer f.Close()
	}

	invalids, halted, err := csvlint.Validate(f, comma, *lazyquotes)
	if err != nil {
		panic(err)
	}
	if len(invalids) == 0 {
		if *verbose {
			fmt.Println("file is valid")
		}
		os.Exit(0)
	}
	for _, invalid := range invalids {
		fmt.Fprintf(os.Stderr, "%s:%d:%s\n", flag.Arg(0), invalid.Num, invalid.Error())
	}
	if halted {
		fmt.Println("\nunable to parse any further")
		os.Exit(1)
	}
	os.Exit(2)
}
