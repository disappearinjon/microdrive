// Package mdturbo provides the MicroDrive/Turbo partition map format,
// along with serializer and deserializer functions.
//
// The format is AFAIK undocumented, but the CiderPress source at
// https://github.com/fadden/ciderpress/blob/master/diskimg/MicroDrive.cpp
// contains a partial description.
package mdturbo

import (
	"reflect"
	"strconv"
)

// tagToOffset extracts an offset tag from a field and returns it as an
// int
func tagToOffset(field reflect.StructField) int64 {
	rawOffset := field.Tag.Get("offset")
	offset, err := strconv.ParseInt(rawOffset, 0, 16)
	if err != nil {
		return -1
	}
	return offset
}
