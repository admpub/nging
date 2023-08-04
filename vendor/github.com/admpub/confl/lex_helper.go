package confl

import "strings"

func stripSlashes(s string) string {
	var builder strings.Builder
	size := len(s)
	var skip bool
	var skipOnly bool
	for i, ch := range s {
		if skip {
			builder.WriteRune(ch)
			skip = false
			continue
		}
		if skipOnly {
			skipOnly = false
			continue
		}
		if ch == '\\' {
			if i+1 < size {
				switch s[i+1] {
				case '\\':
					skip = true
				case 't':
					builder.WriteRune('\t')
					skipOnly = true
				case 'n':
					builder.WriteRune('\n')
					skipOnly = true
				case 'r':
					builder.WriteRune('\r')
					skipOnly = true
				case 'u', 'U', 'x':
					builder.WriteRune(ch)
				}
			}
			continue
		}
		builder.WriteRune(ch)
	}
	return builder.String()
}

func addSlashes(s string, b ...rune) string {
	var builder strings.Builder
	size := len(s)
	for i, v := range s {
		if v == '\\' {
			start := i + 1
			var ok bool
			if start < size {
				switch s[start] {
				case 'u':
					if start+5 < size {
						for _, r := range s[start+1 : start+5] {
							if !isHexadecimal(r) {
								ok = false
								break
							}
							ok = true
						}
					}
				case 'U':
					if start+9 < size {
						for _, r := range s[start+1 : start+9] {
							if !isHexadecimal(r) {
								ok = false
								break
							}
							ok = true
						}
					}
				case 'x':
					if start+3 < size {
						for _, r := range s[start+1 : start+3] {
							if !isHexadecimal(r) {
								ok = false
								break
							}
							ok = true
						}
					}
				}
			}
			if !ok {
				builder.WriteRune(v)
			}
		}
		builder.WriteRune(v)
	}
	return builder.String()
}
