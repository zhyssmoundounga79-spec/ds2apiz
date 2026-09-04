package tests

import "testing"

func repairInvalidJSONBackslashes(s string) string {
	if !strings.Contains(s, "\\") {
		return s
	}
	var out strings.Builder
	out.Grow(len(s) + 10)
	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		if runes[i] == '\\' {
			if i+1 < len(runes) {
				next := runes[i+1]
				switch next {
				case '"', '\\', '/', 'b', 'f', 'n', 'r', 't':
					out.WriteRune('\\')
					out.WriteRune(next)
					i++
					continue
				case 'u':
					if i+5 < len(runes) {
						isHex := true
						for j := 1; j <= 4; j++ {
							r := runes[i+1+j]
							if (r < '0' || r > '9') && (r < 'a' || r > 'f') && (r < 'A' || r > 'F') {
								isHex = false
								break
							}
						}
						if isHex {
							out.WriteRune('\\')
							out.WriteRune('u')
							for j := 1; j <= 4; j++ {
								out.WriteRune(runes[i+1+j])
							}
							i += 5
							continue
						}
					}
				}
			}
			// Not a valid escape sequence, double it
			out.WriteString("\\\\")
		} else {
			out.WriteRune(runes[i])
		}
	}
	return out.String()
}

func TestRepairInvalidJSONBackslashes(t *testing.T) {
	cases := []struct{ in, want string }{
		{`{"path": "C:\\Users\\name"}`, `{"path": "C:\\Users\\name"}`},
		{`{"cmd": "cd D:\\git_codes"}`, `{"cmd": "cd D:\\git_codes"}`},
		{`{"text": "line1\nline2"}`, `{"text": "line1\nline2"}`},
		{`{"path": "D:\\\\back\\\\slash"}`, `{"path": "D:\\\\back\\\\slash"}`},
		{`{"unicode": "\\u2705"}`, `{"unicode": "\\u2705"}`},
		{`{"invalid_u": "\\u123"}`, `{"invalid_u": "\\\\u123"}`},
	}
	for i, c := range cases {
		if got := repairInvalidJSONBackslashes(c.in); got != c.want {
			t.Fatalf("case %d: input=%q got=%q want=%q", i, c.in, got, c.want)
		}
	}
}
