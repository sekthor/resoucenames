package resourcenames

import "strings"

type NamePattern struct {
	segments []nameSegment
}

type nameSegment struct {
	isParam bool
	value   string
}

func split(s string) []string {
	s = strings.Trim(s, "/")
	return strings.Split(s, "/")
}

func isPlaceholder(s string) bool {
	return len(s) >= 3 && s[0] == '{' && s[len(s)-1] == '}'
}

func FromPattern(pattern string) NamePattern {
	segments := split(pattern)
	out := make([]nameSegment, 0, len(segments))

	for _, seg := range segments {
		if isPlaceholder(seg) {
			out = append(out, nameSegment{
				isParam: true,
				value:   seg[1 : len(seg)-1],
			})
			continue
		}

		out = append(out, nameSegment{
			isParam: false,
			value:   seg,
		})
	}

	return NamePattern{
		segments: out,
	}
}
