package mhist

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/alexmorten/mhist/models"
	"github.com/alexmorten/mhist/proto"
	"google.golang.org/grpc"
)

// ErrMeasurementMissingType is returned when one of the Store endpoints is called without a necessary type
var ErrMeasurementMissingType = errors.New("measurement is not categorical or numerical")

// GrpcHandler handles the grpc endpoints for the MhistServer interface
type GrpcHandler struct {
	server     *Server
	port       int
	grpcServer *grpc.Server

	subs *grpcSubscribers
}

// NewGrpcHandler returns a fully initialized GrpcHandler
func NewGrpcHandler(server *Server, port int) *GrpcHandler {
	return &GrpcHandler{
		server: server,
		port:   port,
		subs:   newGrpcSubscribers(),
	}
}

// Run listens on the given port and handles grpc calls
func (h *GrpcHandler) Run() {
	lisAddr := fmt.Sprintf(":%v", h.port)
	log.Println("grpc_handler running on ", lisAddr)
	lis, err := net.Listen("tcp", lisAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	h.grpcServer = grpc.NewServer()

	proto.RegisterMhistServer(h.grpcServer, h)
	if err := h.grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Notify for the Subscriber interface
// notifies are relayed to the grpc connections that are subscribed to measurements
func (h *GrpcHandler) Notify(name string, measurement models.Measurement) {
	h.subs.forEach(func(s *grpcSubscriber) {
		s.Notify(name, measurement)
	})
}

// Shutdown the GrpcHandler
func (h *GrpcHandler) Shutdown() {
	if h.grpcServer == nil {
		return
	}

	h.grpcServer.Stop()
}

// Store the given measurement in mhist
func (h *GrpcHandler) Store(_ context.Context, message *proto.MeasurementMessage) (*proto.Nothing, error) {
	err := h.handleNewMessage(message)
	if err != nil {
		return nil, err
	}

	return &proto.Nothing{}, nil
}

// StoreStream 'ed measurements in mhist
func (h *GrpcHandler) StoreStream(stream proto.Mhist_StoreStreamServer) error {
	for {
		m, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			log.Println(err)
			return err
		}

		err = h.handleNewMessage(m)
		if err != nil {
			return err
		}
	}
}

// Retrieve the requested measurements
func (h *GrpcHandler) Retrieve(_ context.Context, request *proto.RetrieveRequest) (*proto.RetrieveResponse, error) {
	filterDefinition := models.FilterDefinition{}

	if request.Filter != nil {
		filterDefinition.Names = request.Filter.Names
		filterDefinition.Granularity = time.Duration(request.Filter.GranularityNanos)
	}

	endTs := request.End
	if endTs == 0 {
		endTs = time.Now().UnixNano()
	}

	startTs := request.Start
	if startTs == 0 {
		startTs = endTs - (1 * time.Hour).Nanoseconds()
	}
	responseMap := h.server.store.GetMeasurementsInTimeRange(startTs, endTs, filterDefinition)

	return proto.RetrieveResponseFromMeasurementMap(responseMap), nil
}

// Subscribe to measurements
func (h *GrpcHandler) Subscribe(protoFilter *proto.Filter, stream proto.Mhist_SubscribeServer) error {
	subscription := h.subs.newSubscriber()
	filter := models.NewFilterCollection(protoFilter.ToModel())

	for m := range subscription.notifyChan {
		if !filter.Passes(m.name, m.measurement) {
			continue
		}

		pm := proto.MeasurementFromModel(m.measurement)
		message := &proto.MeasurementMessage{
			Name:        m.name,
			Measurement: pm,
		}
		err := stream.SendMsg(message)
		if err != nil {
			log.Println(err)
			log.Println("removing subscription")

			h.subs.removeSubscriber(subscription)

			return err
		}
	}

	return nil
}

func (h *GrpcHandler) handleNewMessage(message *proto.MeasurementMessage) error {
	m := message.Measurement.ToModelWithDefinedTs()

	if m == nil {
		return ErrMeasurementMissingType
	}
	h.server.store.Add(message.Name, m)
	return nil
}
