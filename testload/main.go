package main

import (
	"context"
	"math/rand"

	"github.com/alexmorten/mhist/models"
	"github.com/alexmorten/mhist/proto"
	"google.golang.org/grpc"
)

var raw []byte = []byte("abcdefghijklmnopqrstuvwxyz")

func main() {
	conn, err := grpc.Dial("localhost:6666", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	c := proto.NewMhistClient(conn)
	stream, err := c.StoreStream(context.Background())
	if err != nil {
		panic(err)
	}
	for {
		var m *proto.MeasurementMessage
		i := rand.Intn(100)
		if i > 30 {
			m = &proto.MeasurementMessage{
				Name:        "test_numerical",
				Measurement: proto.MeasurementFromModel(&models.Numerical{Value: rand.Float64()}),
			}
		} else {
			rand.Shuffle(len(raw), func(i, j int) {
				raw[i], raw[j] = raw[j], raw[i]
			})
			m = &proto.MeasurementMessage{
				Name:        "test_raw",
				Measurement: proto.MeasurementFromModel(&models.Raw{Value: raw}),
			}
		}
		err := stream.Send(m)
		if err != nil {
			panic(err)
		}
	}
}
