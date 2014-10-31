# csvlint

`csvlint` is a library and command-line utility for linting CSV files according to [RFC 4180](http://tools.ietf.org/html/rfc4180).

It assumes that your CSV file has an initial header row.

Everything in this README file refers to the command-line utility.
For information about the library, see [godoc](http://godoc.org/github.com/Clever/csvlint).

## Installing

Standalone executables for multiple platforms are available via [Github Releases](https://github.com/Clever/csvlint/releases).

You can also compile from source:

```shell
go get github.com/Clever/csvlint/cmd/csvlint
```

## Usage

`csvlint [options] /path/to/csv/file`

### Options

  * delimiter: the field delimiter to default with
    * default: comma
    * valid options: comma, tab
    * if you want anything else, you're probably doing CSVs wrong
  * lazyquotes: allow a quote to appear in an unquoted field and a non-doubled quote to appear in a quoted field. _WARNING: your file may pass linting, but not parse in the way you would expect_

### Examples

```shell
$ csvlint bad_quote.csv
Record #1 has error: bare " in non-quoted-field

unable to parse any further

$ csvlint --lazyquotes bad_quote.csv
file is valid

$ csvlint mult_long_columns.csv
Record #2 has error: wrong number of fields in line
Record #4 has error: wrong number of fields in line

$ csvlint --delimiter=tab mult_long_columns_tabs.csv
Record #2 has error: wrong number of fields in line
Record #4 has error: wrong number of fields in line

$ csvlint one_long_column.csv
Record #2 has error: wrong number of fields in line

$ csvlint perfect.csv
file is valid
```
