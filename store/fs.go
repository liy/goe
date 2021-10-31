package store

import (
	"net/http"
	"os"
)

/*
Simple storage using OS file system
*/
type Simple struct {}

func (s Simple) Open(name string) (http.File, error) {
	return os.Open(name)
}
