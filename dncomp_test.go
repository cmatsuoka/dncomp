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

func TestRFC(t *testing.T) {
	input := []byte{1, 'F', 3, 'I', 'S', 'I', 4, 'A', 'R', 'P', 'A', 0,
		3, 'F', 'O', 'O', 0xc0, 0, 0xc0, 6, 0, 0}

	expect := []string{"F.ISI.ARPA.", "FOO.F.ISI.ARPA.", "ARPA.", "."}

	res, err := Decode(input)

	if err != nil {
		t.Error("error:", error.Error)
	} else if !equal(res, expect) {
		t.Error("decoding error, expect", expect, ", got", res)
	}
}

func TestEmpty(t *testing.T) {
	input := []byte{}

	expect := []string{}

	res, err := Decode(input)

	if err != nil {
		t.Error("error:", error.Error)
	} else if !equal(res, expect) {
		t.Error("decoding error, expect", expect, ", got", res)
	}
}

func TestRoot(t *testing.T) {
	input := []byte{0}

	expect := []string{"."}

	res, err := Decode(input)

	if err != nil {
		t.Error("error:", error.Error)
	} else if !equal(res, expect) {
		t.Error("decoding error, expect", expect, ", got", res)
	}
}

func TestSimple(t *testing.T) {
	input := []byte{1, 'A'}

	expect := []string{"A."}

	res, err := Decode(input)

	if err != nil {
		t.Error("error:", error.Error)
	} else if !equal(res, expect) {
		t.Error("decoding error, expect", expect, ", got", res)
	}
}
