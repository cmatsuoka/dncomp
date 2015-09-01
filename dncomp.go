/*
Package dncomp implements domain name compression according to RFC 1035
section 4.1.4.
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

	// check if pointer
	if data[lstart]&0xc0 == 0xc0 {
		offset := int(data[lstart]&0x3f)<<8 | int(data[lstart+1])
		if offset >= dstart || addLabel(s, data, dstart, offset) < 0 {
			return -1
		}
		return lstart + 2
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
