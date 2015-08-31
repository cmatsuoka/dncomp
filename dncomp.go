// Package dncomp implements domain name compression according to RFC 1035
// section 4.1.4
package dncomp

/*
func Encode(d []string)[]byte {
}
*/

func add(s *string, data []byte, start int) int {
	// check if pointer
	if data[start]&0xc0 == 0xc0 {
		offset := int(data[start]&0x3f)<<8 | int(data[start+1])
		add(s, data, offset)
		return start + 2
	}

	size := int(data[start])
	start++
	end := start + size
	*s += string(data[start:end]) + "."

	if end >= len(data) || data[end] == 0 {
		return end + 1
	}
	return add(s, data, end)
}

func Decode(data []byte) ([]string, error) {
	var s []string
	for i, num := 0, 0; ; {
		if i >= len(data) {
			break
		}

		s = append(s, "")
		i = add(&s[num], data, i)
		num++
	}

	return s, nil
}
