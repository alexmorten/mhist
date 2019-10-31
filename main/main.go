package main

import (
	"flag"

	_ "net/http/pprof" //pprof for performance analysis

	"github.com/alexmorten/mhist"
)

func main() {
	config := mhist.ServerConfig{}
	flag.IntVar(&config.GrpcPort, "grpc_port", 6666, "defines the port on which the grpc handler operates")
	flag.IntVar(&config.DebugPort, "debug_port", 6667, "defines the port on which the debug handler operates")
	flag.IntVar(&config.MemorySize, "memory_size", 64*1024*1024, "defines the amount of memory the memory store limits itself to. Keep in mind that especially GET request can spike the actual memory usage of the process")
	flag.IntVar(&config.DiskSize, "disk_size", 256*1024*1024, "defines the amount of disk space mhist should occupy")

	flag.Parse()
	server := mhist.NewServer(config)
	server.Run()
}
