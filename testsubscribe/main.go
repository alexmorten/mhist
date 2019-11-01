package main

import (
	"context"
	"log"
	"time"

	"github.com/alexmorten/mhist/proto"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:6666", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	c := proto.NewMhistClient(conn)
	stream, err := c.Subscribe(context.Background(), &proto.Filter{GranularityNanos: int64(time.Second)})
	if err != nil {
		panic(err)
	}
	for {
		m, err := stream.Recv()
		if err != nil {
			panic(err)
		}

		log.Println(m.Measurement.ToModelWithDefinedTs())
	}
}
