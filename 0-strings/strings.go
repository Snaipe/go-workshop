package strings

func IndexByte(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}

func Index(s, substr string) int {
	if substr == "" {
		return 0
	}
	for len(s) > 0 {
		i := IndexByte(s, substr[0])
		if i == -1 {
			return -1
		}
		s = s[i:]

		l := min(len(substr), len(s))
		if s[:l] == substr {
			return i
		}
		s = s[1:]
	}
	return -1
}

func Cut(s, sep string) (before, after string, found bool) {
	i := Index(s, sep)
	if i == -1 {
		return s, "", false
	}
	return s[:i], s[i+len(sep):], true
}

func Split(s, sep string) []string {
	var result []string
	for len(s) > 0 {
		before, after, found := Cut(s, sep)
		if sep == "" {
			before, after = after[:1], after[1:]
		}

		s = after
		result = append(result, before)
		if !found {
			break
		}
	}
	return result
}
