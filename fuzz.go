// +build gofuzz

package simpleuuid

import (
	"encoding/json"
	"net/url"
	"time"
)

func FuzzTestNewBytesURLString(data []byte) int {
	_, err := url.Parse(string(data))
	if err != nil {
		return -1
	}
	idA, err := NewBytes(data)
	if err != nil {
		panic(err)
	}

	if idA.Variant() != 0x4 {
		panic("Variant should be 0x4")
	}

	time, err := NewTime(idA.Time())
	if err != nil {
		panic(err)
	}

	if time.Variant() != 0x4 {
		panic("Variant should be 0x")
	}

	return 1
}

func FuzzNewTime(data []byte) int {
	var t time.Time
	if err := t.UnmarshalText(data); err != nil {
		return -1
	}

	idA, err := NewTimeBytes(t.UTC(), data)
	if err != nil {
		return 0
	}
	if idA.Time().UTC().UnixNano() != t.UTC().UnixNano() {
		panic("mismatched times")
	}

	idB := Copy(idA)
	if idA.Compare(idB) != 0 {
		panic("inequivalent copy")
	}

	idA.Time()
	idA.Version()
	idA.Variant()
	idA.OrderedString()
	idA.AlreadyOrderedString()
	idA.Bytes()

	idB = Copy(idA)
	if idA.Compare(idB) != 0 {
		panic("inequivalent copy")
	}
	if idB.String() != idA.String() {
		panic("string mismatch")
	}

	return 1
}

func FuzzNewString(data []byte) int {
	s := string(data)

	idA, err := NewString(s)
	if err != nil {
		panic(err)
	}

	if idA.String() != string(data) {
		panic("inequivalent string")
	}

	idA.Time()
	idA.Version()
	idA.Variant()
	idA.OrderedString()
	idA.OrderedBytes()
	idA.AlreadyOrderedString()
	idA.Bytes()

	idB := Copy(idA)
	if idA.Compare(idB) != 0 {
		panic("inequivalent copy")
	}
	if idB.String() != idA.String() {
		panic("string mismatch")
	}

	return 1
}

func FuzzNewBytes(data []byte) int {
	idA, err := NewBytes(data)
	if err != nil {
		panic(err)
	}

	idA.Time()
	idA.Version()
	idA.Variant()
	idA.OrderedString()
	idA.OrderedBytes()
	idA.AlreadyOrderedString()
	idAbs := idA.Bytes()

	bs := []byte{}
	err = json.Unmarshal(bs, &idA)
	if err != nil {
		panic(err)
	}
	if len(bs) != len(idAbs) {
		panic("mismatched byte length")
	}
	for string(bs) != string(idAbs) {
		panic("mismatched bytes")
	}

	idB := Copy(idA)
	if idA.Compare(idB) != 0 {
		panic("inequivalent copy")
	}
	if idB.String() != idA.String() {
		panic("string mismatch")
	}

	return 1
}
