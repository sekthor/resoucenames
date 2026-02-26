package resourcenames

import "errors"

var (
	ErrSegmentLengthMismatch   = errors.New("resource name segment count does not match that of the pattern")
	ErrSegmentConstantMismatch = errors.New("a non-variable segment did not match constant pattern segment")
	ErrNotAStruct              = errors.New("resource must be a pointer to a struct")
	ErrMissingSegment          = errors.New("resource is missing variable segment")
)
