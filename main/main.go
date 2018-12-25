package main

import (
	"flag"
	"strings"

	_ "net/http/pprof" //pprof for performance analysis

	"github.com/alexmorten/mhist"
)

func main() {
	config := mhist.ServerConfig{}
	replicationConfigString := ""
	flag.IntVar(&config.HTTPPort, "http_port", 6666, "defines the port on which the http handler operates")
	flag.IntVar(&config.TCPPort, "tcp_port", 6667, "defines the port on which the tcp handler operates")
	flag.IntVar(&config.MemorySize, "memory_size", 64*1024*1024, "defines the amount of memory the memory store limits itself to. Keep in mind that especially GET request can spike the actual memory usage of the process")
	flag.IntVar(&config.DiskSize, "disk_size", 256*1024*1024, "defines the amount of disk space mhist should occupy")
	flag.StringVar(&replicationConfigString, "replicate_to", "", "defines the addresses to replicate to, comma seperated")

	flag.Parse()
	if replicationConfigString != "" {
		config.ReplicationAddresses = strings.Split(replicationConfigString, ",")
	}
	server := mhist.NewServer(config)
	server.Run()
}
