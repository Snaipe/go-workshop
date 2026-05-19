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
	panic("unimplemented")
}

func Cut(s, sep string) (before, after string, found bool) {
	panic("unimplemented")
}

func Split(s, sep string) []string {
	panic("unimplemented")
}
