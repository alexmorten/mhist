package main

import (
	"github.com/codeuniversity/ppp-mhist"
)

func main() {
	server := mhist.NewServer()
	server.Run()
}
