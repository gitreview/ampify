// Public Domain (-) 2010-2011 The Ampify Authors.
// See the Ampify UNLICENSE file for details.

package argo

import (
	"amp/big"
	"bytes"
	"fmt"
	"testing"
)

func Buffer() *bytes.Buffer {
	return bytes.NewBuffer([]byte{})
}

func TestWriteVarint(t *testing.T) {

	tests := map[uint64]string{
		0:                    "\x00",
		123456789:            "\x95\x9a\xef:",
		18446744073709551615: "\xff\xff\xff\xff\xff\xff\xff\xff\xff\x01",
	}

	for value, expected := range tests {
		buf := &bytes.Buffer{}
		enc := &Encoder{buf}
		enc.WriteVarint(value)
		if string(buf.Bytes()) != expected {
			t.Errorf("Got unexpected encoding for %d: %q", value, buf.Bytes())
		}
	}

}

func TestStringArray(t *testing.T) {

	input := []string{"hello", "world", "hehe", "okay"}
	buf := &bytes.Buffer{}
	enc := &Encoder{buf}
	dec := &Decoder{buf}

	err := enc.WriteStringArray(input)
	if err != nil {
		t.Errorf("Got error encoding string array: %s", err)
		return
	}

	result, err := dec.ReadStringArray()
	if err != nil {
		t.Errorf("Got error decoding string array: %s", err)
		return
	}

	if len(input) != len(result) {
		t.Errorf("Got mis-matched result for string array: %#v -> %#v", input, result)
		return
	}

	for idx, item := range input {
		if item != result[idx] {
			t.Errorf("Got mis-matched result for string array: %#v -> %#v", input, result)
			return
		}
	}

}

func TestReadVarint(t *testing.T) {

	tests := map[string]uint64{
		"\x00":                                     0,
		"\x95\x9a\xef:":                            123456789,
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\x01": 18446744073709551615,
	}

	for value, expected := range tests {
		buf := bytes.NewBuffer([]byte(value))
		dec := &Decoder{buf}
		result, err := dec.ReadVarint()
		if err != nil {
			t.Errorf("Got error decoding %q: %s", value, err)
		}
		if result != expected {
			t.Errorf("Got unexpected decoding for %q: %d", value, result)
		}
	}

}

func testWriteInt(t *testing.T) {
	N := int64(8322944)
	buf := Buffer()
	WriteInt(N, buf)
	fmt.Printf("%q\n", string(buf.Bytes()))
}

func testWriteBigInt(t *testing.T) {
	N := big.NewInt(8322944)
	buf := Buffer()
	WriteBigInt(N, buf)
	fmt.Printf("%q\n", string(buf.Bytes()))
}


func testWriteIntOrdering(t *testing.T) {

	buf := Buffer()
	WriteInt(-10258176, buf)
	prev := string(buf.Bytes())

	var i int64

	for i = -10258175; i < 10258175; i++ {
		buf.Reset()
		WriteInt(i, buf)
		cur := string(buf.Bytes())
		if prev >= cur {
			t.Errorf("Lexicographical ordering failure for %d -- %q >= %q", i, prev, cur)
		}
		prev = cur
	}

}

func testWriteBigIntOrdering(t *testing.T) {

	buf := Buffer()
	WriteBigInt(big.NewInt(-10258176), buf)
	prev := string(buf.Bytes())

	var i int64

	for i = -10258175; i < 10258175; i++ {
		buf.Reset()
		WriteBigInt(big.NewInt(i), buf)
		cur := string(buf.Bytes())
		if prev >= cur {
			t.Errorf("Lexicographical ordering failure for %d -- %q >= %q", i, prev, cur)
		}
		prev = cur
	}

}

func decimal(value string) *big.Decimal {
	dec, _ := big.NewDecimal(value)
	return dec
}

func TestWriteDecimalOrdering(t *testing.T) {

	buf := Buffer()
	WriteDecimal(decimal("0"), buf)
	prev := string(buf.Bytes())

	tests := []string{
		"0.02",
		"0.0201",
		"0.05",
		"2",
		"2.30001",
		"2.30002",
	}

	for _, value := range tests {
		buf.Reset()
		WriteDecimal(decimal(value), buf)
		cur := string(buf.Bytes())
		if prev >= cur {
			left, right := decimal(value).Components()
			t.Errorf("Lexicographical ordering failure for %s (%s, %s) -- %q >= %q",
				value, left, right, prev, cur)
		}
		prev = cur
	}

}

func BenchmarkWriteVarint(b *testing.B) {
	buf := Buffer()
	enc := &Encoder{buf}
	for i := 0; i < b.N; i++ {
		buf.Reset()
		enc.WriteVarint(123456789)
	}
}

func BenchmarkWriteNumber(b *testing.B) {
	buf := Buffer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		WriteNumber("123456789", buf)
	}
}
