package main

import (
	_ "net/http/pprof"

	_ "github.com/pkg/profile"
)

func main() {
	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:3000", nil))
	// }()

	startService()
}