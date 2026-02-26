package resourcenames

import "strings"

// Represents all segments of a resource name pattern.
// It is a template to parse resource names against and extract
// the values of variable segments.
type NamePattern struct {
	segments []nameSegment
}

// A segment is a part of a resource name delimited by a forward slash.
// The segment of a resource name pattern can be a constant (e.g. `/resources`)
// or a variable (`{resource_id}`).
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

// FromPattern constructs a NamePattern from a given resource name pattern string.
// Given a pattern string like `/resource/{resource_id}`, it constructs a NamePattern
// consisting of the constant segment `resource` and the variable segment `resource_id`
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
