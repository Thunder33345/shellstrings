package shellstrings

func Parse(input string) []string {
	var parsed []string
	var buffer = ""
	var escaped, doubleQuoted, capture bool
	var length = len(input)
	for i, r := range input {

		switch r {
		case '\\':
			if escaped {
				buffer += "\\\\"
				escaped = false
			} else if isEnd(i, length) {
				buffer += string(r)
				continue
			} else {
				escaped = true
			}
			continue
		case '"':
			if !doubleQuoted && isEnd(i, length) {
				buffer += string(r)
				continue
			}
			if escaped {
				buffer += string(r)
				escaped = false
				continue
			}
			if !doubleQuoted && (seekFor(input[(i+1):], `"`) <= -1) {
				buffer += string(r)
				doubleQuoted = false
				continue
			}

			doubleQuoted = !doubleQuoted
			continue
		case ' ', '\t', '\r', '\n':
			if doubleQuoted {
				buffer += string(r)
				continue
			} else {
				capture = true
			}
		}

		if escaped {
			buffer += "\\" + string(r)
			escaped = false
		} else if capture {
			if len(buffer) > 0 {
				parsed = append(parsed, buffer)
				capture = false
				buffer = ""
			} else {
				capture = false
			}
		} else {
			buffer += string(r)
		}

	}
	if buffer != "" {
		parsed = append(parsed, buffer)
	}
	return parsed
}

func isEnd(pos, total int) bool {
	return pos == (total - 1)
}

func seekFor(stack string, needle string) int {
	//todo better method
	//todo ignore \\"
	for i, c := range stack {
		if string(c) == needle {
			if i > 1 && string(stack[i-1]) == `\` {
				//ignore if it's escaped
				//but return if it's escaped escaped
				if i > 2 && string(stack[i-2]) == `\` {
					return i
				}
				continue
			} else {
				return i
			}
		}
	}
	return -1
}
