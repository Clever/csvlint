package main

import (
	"flag"
	"fmt"
	"github.com/Clever/csvlint"
	"os"
	"strconv"
	"unicode/utf8"
)

func printHelpAndExit(code int) {
	flag.PrintDefaults()
	os.Exit(code)
}

func main() {
	delimiter := flag.String("delimiter", ",", "field delimiter in the file, for instance '\\t' or '|'")
	lazyquotes := flag.Bool("lazyquotes", false, "try to parse improperly escaped quotes")
	help := flag.Bool("help", false, "print help and exit")
	flag.Parse()

	if *help {
		printHelpAndExit(0)
	}

	converted_delimiter, err := strconv.Unquote(`'` + *delimiter + `'`)
	if err != nil {
		fmt.Printf("error unquoting delimiter '%s', note that only one-character delimiters are supported\n\n", *delimiter)
		printHelpAndExit(1)
	}
	// don't need to check size since Unquote returns one-character string
	comma, _ := utf8.DecodeRuneInString(converted_delimiter)

	if len(flag.Args()) != 1 {
		fmt.Println("csvlint accepts a single filepath as an argument\n")
		printHelpAndExit(1)
	}

	f, err := os.Open(flag.Args()[0])
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("file '%s' does not exist\n", flag.Args()[0])
			os.Exit(1)
		} else {
			panic(err)
		}
	}
	defer f.Close()

	invalids, halted, err := csvlint.Validate(f, comma, *lazyquotes)
	if err != nil {
		panic(err)
	}
	if len(invalids) == 0 {
		fmt.Println("file is valid")
		os.Exit(0)
	}
	for _, invalid := range invalids {
		fmt.Println(invalid.Error())
	}
	if halted {
		fmt.Println("\nunable to parse any further")
		os.Exit(1)
	}
	os.Exit(2)
}
