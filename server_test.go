package mhist

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/alexmorten/mhist/models"
	"github.com/alexmorten/mhist/proto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Server(t *testing.T) {
	formerDataPath := dataPath
	dataPath = "test_data"
	defer func() {
		os.RemoveAll(dataPath)
		dataPath = formerDataPath
	}()
	server := NewServer(ServerConfig{MemorySize: 2 * 1024, DiskSize: 24 * 1024 * 1024})

	t.Run("GrpcHandler", func(t *testing.T) {
		t.Run("Storing and Retrieving measurements without filters", func(t *testing.T) {
			numericalValues := []float64{60, 10, 40, 20, 50, 42}
			categoricalValues := []string{"a", "b", "a", "de", "c", "b", "a"}
			rawValues := [][]byte{[]byte("some_raw_value idk"), []byte("some_raw_value i still dont know"), []byte("some"), []byte("thing")}
			serverTestSetup(t, server, numericalValues, categoricalValues, rawValues, func(_ int) int64 { return 0 })
			server.store.diskStore.commit()
			request := &proto.RetrieveRequest{}
			response, err := server.grpcHandler.Retrieve(context.Background(), request)
			require.NoError(t, err)
			require.NotNil(t, response)
			numericalResponseValues := []float64{}
			hist := response.Histories["some_name"]
			require.NotNil(t, hist)
			for _, measurement := range hist.Measurements {
				numericalResponseValues = append(numericalResponseValues, measurement.Type.(*proto.Measurement_Numerical).Numerical.Value)
			}

			assert.ElementsMatch(t, numericalResponseValues, numericalValues)

			categoricalResponseValues := []string{}
			for _, measurement := range response.Histories["some_other_name"].Measurements {
				categoricalResponseValues = append(categoricalResponseValues, measurement.Type.(*proto.Measurement_Categorical).Categorical.Value)
			}

			assert.ElementsMatch(t, categoricalResponseValues, categoricalValues)

			rawResponseValues := [][]byte{}
			for _, measurement := range response.Histories["some_even_different_name"].Measurements {
				rawResponseValues = append(rawResponseValues, measurement.Type.(*proto.Measurement_Raw).Raw.Value)
			}

			assert.ElementsMatch(t, rawResponseValues, rawValues)
		})
	})
}

func Test_ServerFilter(t *testing.T) {
	formerDataPath := dataPath
	dataPath = "test_data"
	defer func() {
		os.RemoveAll(dataPath)
		dataPath = formerDataPath
	}()
	server := NewServer(ServerConfig{MemorySize: 24 * 1024 * 1024, DiskSize: 24 * 1024 * 1024})

	t.Run("GrpcHandler", func(t *testing.T) {
		t.Run("Storing and Retrieving measurements with filters", func(t *testing.T) {
			numericalValues := []float64{60, 10, 40, 20, 50, 42}
			categoricalValues := []string{"a", "b", "a", "de", "c", "b", "a"}
			rawValues := [][]byte{[]byte("some_raw_value idk"), []byte("some_raw_value i still dont know"), []byte("some"), []byte("thing")}

			serverTestSetup(t, server, numericalValues, categoricalValues, rawValues, func(i int) int64 { return int64(i*1000 + 1000) })
			request := &proto.RetrieveRequest{
				Start: 2001,
			}
			response, err := server.grpcHandler.Retrieve(context.Background(), request)
			require.NoError(t, err)

			numericalResponseValues := []float64{}

			for _, measurement := range response.Histories["some_name"].Measurements {
				numericalResponseValues = append(numericalResponseValues, measurement.Type.(*proto.Measurement_Numerical).Numerical.Value)
			}

			assert.ElementsMatch(t, numericalResponseValues, numericalValues[2:])

			categoricalResponseValues := []string{}
			for _, measurement := range response.Histories["some_other_name"].Measurements {
				categoricalResponseValues = append(categoricalResponseValues, measurement.Type.(*proto.Measurement_Categorical).Categorical.Value)
			}

			assert.ElementsMatch(t, categoricalResponseValues, categoricalValues[2:])

			rawResponseValues := [][]byte{}

			for _, measurement := range response.Histories["some_even_different_name"].Measurements {
				rawResponseValues = append(rawResponseValues, measurement.Type.(*proto.Measurement_Raw).Raw.Value)
			}

			assert.ElementsMatch(t, rawResponseValues, rawValues[2:])
		})
	})
}

func Test_ServerFineWithMultipleCommits(t *testing.T) {
	formerDataPath := dataPath
	dataPath = "test_data"
	defer func() {
		os.RemoveAll(dataPath)
		dataPath = formerDataPath
	}()
	server := NewServer(ServerConfig{MemorySize: 64 * 1024, DiskSize: 24 * 1024 * 1024})

	t.Run("GrpcHandler", func(t *testing.T) {
		t.Run("Storing and Retrieving measurements with filters", func(t *testing.T) {
			numericalValues := []float64{60, 10, 40, 20, 50, 42}
			categoricalValues := []string{"a", "b", "a", "de", "c", "b", "a"}
			rawValues := [][]byte{[]byte("some_raw_value idk"), []byte("some_raw_value i still dont know"), []byte("some"), []byte("thing")}

			for i := 0; i < 100000; i++ {
				serverTestSetup(t, server, numericalValues, categoricalValues, rawValues, func(i int) int64 { return 0 })
			}
			log.Println("reading")
			request := &proto.RetrieveRequest{}
			response, err := server.grpcHandler.Retrieve(context.Background(), request)
			log.Println("reading done")
			require.NoError(t, err)
			require.NotNil(t, response.Histories["some_name"])
			require.NotNil(t, response.Histories["some_other_name"])
			require.NotNil(t, response.Histories["some_even_different_name"])
			assert.Len(t, response.Histories["some_name"].Measurements, 600000)
			assert.Len(t, response.Histories["some_other_name"].Measurements, 700000)
			assert.Len(t, response.Histories["some_even_different_name"].Measurements, 400000)
		})
	})
}

func serverTestSetup(t *testing.T, server *Server, numericalValues []float64, categoricalValues []string, rawValues [][]byte, tsForIndex func(i int) int64) {
	for i, value := range numericalValues {
		message := &proto.MeasurementMessage{
			Name:        "some_name",
			Measurement: proto.MeasurementFromModel(&models.Numerical{Ts: tsForIndex(i), Value: value}),
		}
		_, err := server.grpcHandler.Store(context.Background(), message)
		assert.NoError(t, err)

	}

	for i, value := range categoricalValues {
		message := &proto.MeasurementMessage{
			Name:        "some_other_name",
			Measurement: proto.MeasurementFromModel(&models.Categorical{Ts: tsForIndex(i), Value: value}),
		}
		_, err := server.grpcHandler.Store(context.Background(), message)
		assert.NoError(t, err)

	}

	for i, value := range rawValues {
		message := &proto.MeasurementMessage{
			Name:        "some_even_different_name",
			Measurement: proto.MeasurementFromModel(&models.Raw{Ts: tsForIndex(i), Value: value}),
		}
		_, err := server.grpcHandler.Store(context.Background(), message)
		assert.NoError(t, err)

	}
}
