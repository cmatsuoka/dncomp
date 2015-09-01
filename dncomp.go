// Package dncomp implements domain name compression according to RFC 1035
// section 4.1.4
package dncomp

import "errors"

/*
func Encode(d []string)[]byte {
}
*/

type ptrMap map[int]bool

func addLabel(s *string, data []byte, start int, p ptrMap) int {
	dataSize := len(data)

	if start >= dataSize {
		return -1
	}

	// check if pointer
	if data[start]&0xc0 == 0xc0 {
		offset := int(data[start]&0x3f)<<8 | int(data[start+1])

		// loop detection
		if p[offset] {
			return -1
		}
		p[offset] = true

		if offset >= dataSize || addLabel(s, data, offset, p) < 0 {
			return -1
		}
		return start + 2
	}

	labelSize := int(data[start])
	if start + labelSize >= dataSize {
		return -1
	}

	start++
	end := start + labelSize

	*s += string(data[start:end])

	if end >= dataSize || data[end] == 0 {
		return end + 1
	}

	*s += "."

	return addLabel(s, data, end, p)
}

func Decode(data []byte) ([]string, error) {
	var s []string
	
	for i, num := 0, 0; ; {
		if i >= len(data) {
			break
		}

		// use this to prevent pointer loops
		p := make(ptrMap)

		s = append(s, "")
		i = addLabel(&s[num], data, i, p)

		if (i < 0) {
			return nil, errors.New("malformed compressed data")
		}

		num++
	}

	return s, nil
}
