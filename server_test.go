package mhist

import (
	"context"
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
	server := NewServer(ServerConfig{MemorySize: 24 * 1024 * 1024, DiskSize: 24 * 1024 * 1024})

	t.Run("GrpcHandler", func(t *testing.T) {
		t.Run("Storing and Retrieving measurements without filters", func(t *testing.T) {
			numericalValues := []float64{60, 10, 40, 20, 50, 42}
			categoricalValues := []string{"a", "b", "a", "de", "c", "b", "a"}
			serverTestSetup(t, server, numericalValues, categoricalValues, func(_ int) int64 { return 0 })
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
			serverTestSetup(t, server, numericalValues, categoricalValues, func(i int) int64 { return int64(i*1000 + 1000) })
			request := &proto.RetrieveRequest{
				Start: 2001,
			}
			response, err := server.grpcHandler.Retrieve(context.Background(), request)
			require.NoError(t, err)

			numericalResponseValues := []float64{}
			hist := response.Histories["some_name"]
			require.NotNil(t, hist)

			for _, measurement := range hist.Measurements {
				numericalResponseValues = append(numericalResponseValues, measurement.Type.(*proto.Measurement_Numerical).Numerical.Value)
			}

			assert.ElementsMatch(t, numericalResponseValues, numericalValues[2:])

			categoricalResponseValues := []string{}
			for _, measurement := range response.Histories["some_other_name"].Measurements {
				categoricalResponseValues = append(categoricalResponseValues, measurement.Type.(*proto.Measurement_Categorical).Categorical.Value)
			}

			assert.ElementsMatch(t, categoricalResponseValues, categoricalValues[2:])
		})
	})
}

func serverTestSetup(t *testing.T, server *Server, numericalValues []float64, categoricalValues []string, tsForIndex func(i int) int64) {

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
}
