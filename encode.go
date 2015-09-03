package dncomp

import (
	"bytes"
	"errors"
	"strings"
)

func splitDomainName(s string) (string, string) {
	parts := strings.Split(s, ".")

	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], strings.Join(parts[1:], ".")
}

func encodeLabel(b *bytes.Buffer, s string) (int, error) {
	labelSize := len(s)
	if labelSize > 63 {
		return 0, errors.New("label is too long")
	}

	b.WriteByte(byte(labelSize))
	if labelSize == 0 {
		return 1, nil
	}

	_, err := b.WriteString(s)
	if err != nil {
		return 0, errors.New("label encoding error")
	}
	return labelSize + 1, nil
}

func encodeDomainName(b *bytes.Buffer, s string, i int) (int, string, error) {
	var head, tail string

	tail = s
	size := 0
	var err error
	for {
		head, tail = splitDomainName(tail)
		size, err = encodeLabel(b, head)
		if err != nil {
			return 0, s, err
		}

		if head == "" {
			break
		}
	}

	return i + size, tail, nil
}

// Encode receives a list of domain names and encodes it according to
// the compression scheme described in RFC 1035 section 4.1.4.
func Encode(d []string) ([]byte, error) {
	var b bytes.Buffer

	index := 0
	for _, s := range d {
		n, _, err := encodeDomainName(&b, s, index)
		if err != nil {
			return nil, err
		}
		index += n
	}

	return b.Bytes(), nil
}
