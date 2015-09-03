package dncomp

import (
	"bytes"
	"testing"
)

func checkEncode(t *testing.T, input []string, err error, res, expect []byte) {
	switch {
	case err != nil:
		t.Error("error:", err.Error())
	case !bytes.Equal(res, expect):
		t.Error("encoding error for input", input,
			", expect", expect, ", got", res)
	}
}

// Encoding tests

func TestEncodeSimple(t *testing.T) {
	input := []string{"A"}

	expect := []byte{1, 'A', 0}

	res, err := Encode(input)
	checkEncode(t, input, err, res, expect)
}

func TestEncodeLongDomain(t *testing.T) {
	input := []string{"A.B.CITS.BR"}

	expect := []byte{1, 'A', 1, 'B', 4, 'C', 'I', 'T', 'S', 2, 'B', 'R', 0}

	res, err := Encode(input)
	checkEncode(t, input, err, res, expect)
}

func TestEncodeRoot(t *testing.T) {
	input := []string{""}

	expect := []byte{0}

	res, err := Encode(input)
	checkEncode(t, input, err, res, expect)
}

func TestEncodeTwoDomains(t *testing.T) {
	input := []string{"A.B", "C.D"}

	expect := []byte{1, 'A', 1, 'B', 0, 1, 'C', 1, 'D', 0}

	res, err := Encode(input)
	checkEncode(t, input, err, res, expect)
}
