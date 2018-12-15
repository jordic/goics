package goics

// splitLength returns a slice of strings, each string is at most length bytes long.
// Does not break UTF-8 codepoints.
func splitLength(s string, length int) []string {
	var ret []string

	for len(s) > 0 {
		tmp := truncateString(s, length)
		if len(tmp) == 0 {
			// too short length, or invalid UTF-8 string
			break
		}
		ret = append(ret, tmp)
		s = s[len(tmp):]
	}

	return ret
}

// truncateString truncates s to a maximum of length bytes without breaking UTF-8 codepoints.
func truncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}

	// UTF-8 continuation bytes start with 10xx xxxx:
	// 0xc0 = 1100 0000
	// 0x80 = 1000 0000
	cutoff := length
	for s[cutoff]&0xc0 == 0x80 {
		cutoff--
		if cutoff < 0 {
			cutoff = 0
			break
		}
	}

	return s[0:cutoff]
}
