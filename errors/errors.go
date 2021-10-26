package errors

import (
	"errors"
)

var (
	ErrObjectNotFound     = errors.New("object not found")
	ErrInvalidObjectType  = errors.New("invalid object type")
	ErrReferenceNoteFound = errors.New("reference not found")
	// Happens when decode raw object which has wrong type, e.g., try to
	// decode tag raw object into a a commit
	ErrRawObjectTypeWrong = errors.New("raw object type is wrong")
)