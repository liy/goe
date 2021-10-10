package errors

import "errors"

var (
	ErrObjectNotFound = errors.New("object not found")
	ErrInvalidObjectType = errors.New("invalid object type")
	ErrReferenceNoteFound = errors.New("reference not found")
)
