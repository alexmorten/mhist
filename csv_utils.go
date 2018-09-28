package mhist

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const fieldSeperatorRune = ','
const fieldSeperator = string(fieldSeperatorRune)

const fieldSeperatorSize = len(fieldSeperator)

const newLineSize = len("\n")

func constructCsvLine(id int64, m Measurement) ([]byte, error) {
	value := m.ValueString()
	tsString := strconv.FormatInt(m.Timestamp(), 10)
	idString := strconv.FormatInt(id, 10)

	if strings.ContainsRune(value, ',') || strings.ContainsRune(value, '\n') {
		return nil, fmt.Errorf("'%v' contains an invalid char", value)
	}

	byteSize := len(idString) + fieldSeperatorSize + len(tsString) + fieldSeperatorSize + len(value) + newLineSize
	byteSlice := make([]byte, 0, byteSize)

	byteSlice = append(byteSlice, []byte(idString)...)
	byteSlice = append(byteSlice, []byte(fieldSeperator)...)
	byteSlice = append(byteSlice, []byte(tsString)...)
	byteSlice = append(byteSlice, []byte(fieldSeperator)...)
	byteSlice = append(byteSlice, []byte(value)...)
	byteSlice = append(byteSlice, []byte("\n")...)

	return byteSlice, nil
}

func newCsvReader(r io.Reader) *csv.Reader {
	reader := csv.NewReader(r)
	reader.Comma = fieldSeperatorRune
	return reader
}
