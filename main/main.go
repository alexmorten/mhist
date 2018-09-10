package main

import (
	"github.com/codeuniversity/ppp-mhist"
)

const memorySize = 64 * 1024 * 1024 //64MB for now, should be filled by commandline argument

func main() {
	server := mhist.NewServer(memorySize)
	server.Run()
}
