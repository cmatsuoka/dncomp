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

// ErrMalformedCompressedData is returned when compressed domain names
// cannot be properly parsed due to invalid offsets, label sizes, or
// character set.
var ErrMalformedCompressedData = errors.New("dncomp: malformed compressed data")

// decodeLabel adds a label from offset off of compressed data to byte
// buffer b, for the name record starting at start.
func decodeLabel(b *bytes.Buffer, data []byte, start, off int) int {
	for {
		if off >= len(data) {
			return -1
		}

		size := int(data[off])
		off++

		// check for pointer
		switch size & 0xc0 {
		case 0x00:
			// end of domain name
			if size == 0 {
				return off
			}
			// sanity check: invalid label length
			if off+size > len(data) {
				return -1
			}

			// write label to domain name string
			end := off + size
			_, err := b.Write(data[off:end])
			if err != nil {
				return -1
			}
			off = end

			// sanity check: missing end marker
			if end >= len(data) {
				return -1
			}

			// add dot if we have more labels
			if data[end] != 0 {
				b.WriteByte(byte('.'))
			}

		case 0xc0:
			// sanity check: second octet available
			if off >= len(data) {
				return -1
			}

			// pointer offset
			ptr := (size&0x3f)<<8 | int(data[off])
			if ptr >= start || decodeLabel(b, data, start, ptr) < 0 {
				return -1
			}
			return off + 1

		default:
			// reserved codes 10 and 01
			return -1
		}
	}
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
			return nil, ErrMalformedCompressedData
		}

		s = append(s, b.String())
	}

	return s, nil
}
