package tests

import (
	"net/http"
)

type Embeded struct{
}

func (e Embeded) Open(name string) (http.File, error) {
	return FS(false).Open(name)
}