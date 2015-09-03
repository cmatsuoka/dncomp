package dncomp

import "testing"

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func checkDecode(t *testing.T, input []byte, err error, res, expect []string) {
	switch {
	case err != nil:
		t.Error("error:", err.Error())
	case !equal(res, expect):
		t.Error("decoding error for input", input,
			", expect", expect, ", got", res)
	}
}

func checkError(t *testing.T, input []byte, err error, res []string) {
	if err == nil {
		t.Error("error return failed for input", input, ", got", res)
	}
}

// Decoding tests

func TestDecodeRFCExample(t *testing.T) {
	input := []byte{1, 'F', 3, 'I', 'S', 'I', 4, 'A', 'R', 'P', 'A', 0,
		3, 'F', 'O', 'O', 0xc0, 0, 0xc0, 6, 0}

	expect := []string{"F.ISI.ARPA", "FOO.F.ISI.ARPA", "ARPA", ""}

	res, err := Decode(input)
	checkDecode(t, input, err, res, expect)
}

func TestDecodeEmpty(t *testing.T) {
	input := []byte{}

	expect := []string{}

	res, err := Decode(input)
	checkDecode(t, input, err, res, expect)
}

func TestDecodeRoot(t *testing.T) {
	input := []byte{0}

	expect := []string{""}

	res, err := Decode(input)
	checkDecode(t, input, err, res, expect)
}

func TestDecodeSimple(t *testing.T) {
	input := []byte{1, 'A', 0}

	expect := []string{"A"}

	res, err := Decode(input)
	checkDecode(t, input, err, res, expect)
}

func TestDecodeInvalidLabelLength1(t *testing.T) {
	// invalid label length
	input := []byte{1}
	res, err := Decode(input)
	checkError(t, input, err, res)
}

func TestDecodeInvalidLabelLength2(t *testing.T) {
	// invalid label length
	input := []byte{'A'}
	res, err := Decode(input)
	checkError(t, input, err, res)
}

func TestDecodeInvalidLabelLength3(t *testing.T) {
	// invalid label length
	input := []byte{5, 'A', 'B'}
	res, err := Decode(input)
	checkError(t, input, err, res)
}

func TestDecodeInvalidLabelLength4(t *testing.T) {
	// invalid length in second label
	input := []byte{2, 'A', 'B', 0, 1}
	res, err := Decode(input)
	checkError(t, input, err, res)
}

func TestDecodeMissingEnd(t *testing.T) {
	// invalid length in second label
	input := []byte{1, 'A'}
	res, err := Decode(input)
	checkError(t, input, err, res)
}

func TestDecodeInvalidPointer1(t *testing.T) {
	// pointer to invalid offset
	input := []byte{2, 'A', 'B', 0, 0xc0, 6}
	res, err := Decode(input)
	checkError(t, input, err, res)
}

func TestDecodeInvalidPointer2(t *testing.T) {
	// pointer to the pointer offset
	input := []byte{2, 'A', 'B', 0, 0xc0, 5}
	res, err := Decode(input)
	checkError(t, input, err, res)
}

func TestDecodeInvalidPointer3(t *testing.T) {
	// forward loop
	input := []byte{2, 'A', 'B', 0, 0xc0, 6, 'C', 'D'}
	res, err := Decode(input)
	checkError(t, input, err, res)
}

func TestDecodeInvalidPointer4(t *testing.T) {
	// pointer starting with 10
	input := []byte{2, 'A', 'B', 0, 0x80, 0}
	res, err := Decode(input)
	checkError(t, input, err, res)
}

func TestDecodeInvalidPointer5(t *testing.T) {
	// pointer starting with 01
	input := []byte{2, 'A', 'B', 0, 0x40, 0}
	res, err := Decode(input)
	checkError(t, input, err, res)
}

func TestDecodeInvalidPointer6(t *testing.T) {
	// pointer with only one octet
	input := []byte{0xc0}
	res, err := Decode(input)
	checkError(t, input, err, res)
}

func TestDecodeLoop1(t *testing.T) {
	// pointer loop to pointer offset
	input := []byte{2, 'A', 'B', 0, 0xc0, 4}
	res, err := Decode(input)
	checkError(t, input, err, res)
}

func TestDecodeLoop2(t *testing.T) {
	// pointer loop forward and back
	input := []byte{2, 'A', 'B', 0, 0xc0, 7, 'C', 'D', 0xc0, 4}
	res, err := Decode(input)
	checkError(t, input, err, res)
}

func TestDecodeLoop3(t *testing.T) {
	// pointer loop backwards to invalid offset inside the same domain
	input := []byte{2, 'A', 'B', 0xc0, 1}
	res, err := Decode(input)
	checkError(t, input, err, res)
}

func TestDecodeLoop4(t *testing.T) {
	// pointer loop backwards inside the same domain
	input := []byte{2, 'A', 'B', 0xc0, 0}
	res, err := Decode(input)
	checkError(t, input, err, res)
}

func TestDecodeLoop5(t *testing.T) {
	// pointer loop backwards to invalid offset of previous domain
	input := []byte{2, 'A', 'B', 0, 0xc0, 1}
	res, err := Decode(input)
	checkError(t, input, err, res)
}

