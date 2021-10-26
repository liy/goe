package plumbing

import (
	"testing"

	"github.com/liy/goe/tests"
)

func GetRawObject(t *testing.T, hash string) *RawObject {
	raw := NewRawObject(ToHash(hash))
	raw.ReadFile(tests.GetObjectFile(t, hash))
	return raw
}