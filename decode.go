/*
Package dncomp implements domain name compression according to RFC 1035
section 4.1.4.

Compressed domain names are composed by a sequence of label size and labels,
ending with 0 or a pointer to an offset in compressed data. Pointers are
encoded in a two-octet word with the most significant bits set to 1, and the
rest of the word containing the pointer offset. The maximum size of a label
is 63 octets, and pointer codes 10 and 01 are reserved. Pointers should
always point backwards to data that has already been processed.

Example:

	1, 'F',
	3, 'I', 'S', 'I',
	4, 'A', 'R', 'P', 'A', 0,
	3, 'F', 'O', 'O', 0xc0, 0,
	0xc0, 6, 0

decodes to:

	"F.ISI.ARPA"
	"FOO.F.ISI.ARPA"
	"ARPA"
	"" (root domain name, no labels)

*/
package dncomp

import (
	"bytes"
	"errors"
)

// decodeLabel adds a label from offset lstart of compressed data to byte
// buffer b, for the domain name starting at dstart.
func decodeLabel(b *bytes.Buffer, data []byte, dstart, lstart int) int {
	dataSize := len(data)

	if lstart >= dataSize {
		return -1
	}

	// check for pointer
	switch data[lstart] & 0xc0 {
	case 0xc0:
		// check if second octet is available
		if lstart+1 >= dataSize {
			return -1
		}

		// offset
		offset := int(data[lstart]&0x3f)<<8 | int(data[lstart+1])
		if offset >= dstart || decodeLabel(b, data, dstart, offset) < 0 {
			return -1
		}
		return lstart + 2

	case 0x80, 0x40:
		// reserved codes 10 and 01
		return -1
	}

	labelSize := int(data[lstart])

	// check end of domain name
	if labelSize == 0 {
		return lstart + 1
	}

	if lstart+labelSize >= dataSize {
		return -1
	}

	lstart++
	end := lstart + labelSize

	// write label to domain name string
	_, err := b.Write(data[lstart:end])
	if err != nil {
		return -1
	}

	// check for missing end marker
	if end >= dataSize {
		return -1
	}

	// add dot if we have more labels
	if data[end] != 0 {
		b.WriteByte(byte('.'))
	}

	return decodeLabel(b, data, dstart, end)
}

// Decode uncompresses a list of domain names encoded according to RFC 1035
// section 4.1.4, returning a list of uncompressed domain names or an error
// if data is inconsistent.
func Decode(data []byte) ([]string, error) {
	var s []string

	for i := 0; ; {
		if i >= len(data) {
			break
		}

		var b bytes.Buffer
		i = decodeLabel(&b, data, i, i)
		if i < 0 {
			return nil, errors.New("malformed compressed data")
		}

		s = append(s, b.String())
	}

	return s, nil
}
