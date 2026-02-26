package resourcenames

func (p NamePattern) Parse(resourcename string) (map[string]string, error) {
	nameSegments := split(resourcename)
	namedParams := make(map[string]string)

	if len(nameSegments) != len(p.segments) {
		return nil, ErrSegmentLengthMismatch
	}

	for i, patternSegment := range p.segments {
		if !patternSegment.isParam {

			// if we arrive here, the segment is of constant values
			// if the resourcename's segment differs from the pattern segment
			// the patterns don't match
			if patternSegment.value != nameSegments[i] {
				return namedParams, ErrSegmentConstantMismatch
			}

			continue
		}
		namedParams[patternSegment.value] = nameSegments[i]
	}
	return namedParams, nil
}
