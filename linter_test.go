package csvlint

import (
	"encoding/csv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var validationTable = []struct {
	file     string
	err      error
	invalids []CSVError
	comma    rune
	halted   bool
}{
	{file: "./test_data/perfect.csv", err: nil, invalids: []CSVError{}},
	{file: "./test_data/perfect_tab.csv", err: nil, comma: '\t', invalids: []CSVError{}},
	{file: "./test_data/perfect_pipe.csv", err: nil, comma: '|', invalids: []CSVError{}},
	{file: "./test_data/perfect_colon.csv", err: nil, comma: ':', invalids: []CSVError{}},
	{file: "./test_data/perfect_semicolon.csv", err: nil, comma: ';', invalids: []CSVError{}},
	{file: "./test_data/one_long_column.csv", err: nil, invalids: []CSVError{{
		Record: []string{"d", "e", "f", "g"},
		err:    csv.ErrFieldCount,
		Num:    2,
	}}},
	{file: "./test_data/mult_long_columns.csv", err: nil, invalids: []CSVError{
		{
			Record: []string{"d", "e", "f", "g"},
			err:    csv.ErrFieldCount,
			Num:    2,
		}, {
			Record: []string{"k", "l", "m", "n"},
			err:    csv.ErrFieldCount,
			Num:    4,
		}},
	},
	{file: "./test_data/mult_long_columns_tabs.csv", err: nil, comma: '\t', invalids: []CSVError{
		{
			Record: []string{"d", "e", "f", "g"},
			err:    csv.ErrFieldCount,
			Num:    2,
		}, {
			Record: []string{"k", "l", "m", "n"},
			err:    csv.ErrFieldCount,
			Num:    4,
		}},
	},
}

func TestTable(t *testing.T) {
	for _, test := range validationTable {
		f, err := os.Open(test.file)
		assert.Nil(t, err)
		defer f.Close()
		comma := test.comma
		if test.comma == 0 {
			comma = ','
		}
		invalids, halted, err := Validate(f, comma, false)
		assert.Equal(t, test.err, err)
		assert.Equal(t, halted, test.halted)
		assert.Equal(t, test.invalids, invalids)
	}
}

var errTable = []struct {
	err     error
	message string
}{
	{
		err:     CSVError{Record: []string{"a", "b", "c"}, Num: 3, err: csv.ErrFieldCount},
		message: "Record #3 has error: wrong number of fields",
	},
	{
		err:     CSVError{Record: []string{"d", "e", "f"}, Num: 1, err: csv.ErrBareQuote},
		message: `Record #1 has error: bare " in non-quoted-field`,
	},
}

func TestErrors(t *testing.T) {
	for _, test := range errTable {
		assert.Equal(t, test.err.Error(), test.message)
	}
}
