/*
Package dncomp implements domain name compression according to RFC 1035
section 4.1.4.

Compressed domain names are composed by a sequence of label size and labels,
ending with 0 or a pointer to an offset in compressed data. Pointers are
encoded in a two-octet word with the most significant bits set to 1, and the
rest of the word containing the pointer offset. The maximum size of a label
is 63 octets, and pointer codes 10 and 01 are reserved. Pointers should
always point backwards to data that has already been processed.
*/
package dncomp

import "errors"

/*
func Encode(d []string)[]byte {
}
*/

func addLabel(s *string, data []byte, dstart, lstart int) int {
	dataSize := len(data)

	if lstart >= dataSize {
		return -1
	}

	switch data[lstart]&0xc0 {
	case 0xc0:
		// pointer
		offset := int(data[lstart]&0x3f)<<8 | int(data[lstart+1])
		if offset >= dstart || addLabel(s, data, dstart, offset) < 0 {
			return -1
		}
		return lstart + 2

	case 0x80, 0x40:
		// invalid pointer codes 10 and 01
		return -1
	}

	labelSize := int(data[lstart])
	if lstart+labelSize >= dataSize {
		return -1
	}

	lstart++
	end := lstart + labelSize

	*s += string(data[lstart:end])

	if end >= dataSize || data[end] == 0 {
		return end + 1
	}

	*s += "."

	return addLabel(s, data, dstart, end)
}

// Decode uncompresses a list of domain names encoded according to RFC 1035
// section 4.1.4, returning a list of uncompressed domain names or an error
// if data is inconsistent.
func Decode(data []byte) ([]string, error) {
	var s []string

	for i, num := 0, 0; ; {
		if i >= len(data) {
			break
		}

		s = append(s, "")
		i = addLabel(&s[num], data, i, i)

		if i < 0 {
			return nil, errors.New("malformed compressed data")
		}

		num++
	}

	return s, nil
}
