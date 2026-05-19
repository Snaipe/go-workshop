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
	total := 0
	for len(s) > 0 {
		i := IndexByte(s, substr[0])
		if i == -1 {
			break
		}
		s = s[i:]
		total += i

		l := min(len(s), len(substr))
		if s[:l] == substr {
			return total
		}

		s = s[1:]
		total += 1
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

	if sep == "" {
		for i := 0; i < len(s); i++ {
			result = append(result, s[i:i+1])
		}
		return result
	}

	for {
		before, after, found := Cut(s, sep)
		if !found {
			result = append(result, before)
			break
		}

		result = append(result, before)
		s = after
	}
	return result
}

func IndexRune(s string, r rune) int {
	for i, sr := range s {
		if sr == r {
			return i
		}
	}
	return -1
}

func IndexAny(s, chars string) int {
	for i, r := range s {
		if IndexRune(chars, r) != -1 {
			return i
		}
	}
	return -1
}
