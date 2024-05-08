// Package formatter implements some functions to format string, struct.
package formatter

import (
	"encoding/json"
	"io"
	"strconv"
	"strings"

	"github.com/sllt/af/convertor"
	"golang.org/x/exp/constraints"
)

// Comma add comma to a number value by every 3 numbers from right. ahead by symbol char.
// if value is invalid number string eg "aa", return empty string
// Comma("12345", "$") => "$12,345", Comma(12345, "$") => "$12,345"
func Comma[T constraints.Float | constraints.Integer | string](value T, symbol string) string {
	numString := convertor.ToString(value)

	_, err := strconv.ParseFloat(numString, 64)
	if err != nil {
		return ""
	}

	index := strings.Index(numString, ".")
	if index == -1 {
		index = len(numString)
	}

	for index > 3 {
		index = index - 3
		numString = numString[:index] + "," + numString[index:]
	}

	return symbol + numString
}

// Pretty data to JSON string.
func Pretty(v any) (string, error) {
	out, err := json.MarshalIndent(v, "", "    ")
	return string(out), err
}

// PrettyToWriter pretty encode data to writer.
func PrettyToWriter(v any, out io.Writer) error {
	enc := json.NewEncoder(out)
	enc.SetIndent("", "    ")

	if err := enc.Encode(v); err != nil {
		return err
	}

	return nil
}
