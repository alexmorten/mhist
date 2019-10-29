package mhist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

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

	t.Run("HTTPHandler", func(t *testing.T) {
		t.Run("POSTing and GETing measurements works", func(t *testing.T) {
			numericalValues := []float64{10, 60, 40, 20, 50, 42}
			for _, value := range numericalValues {
				w := httptest.NewRecorder()
				body := fmt.Sprintf(`{"name":"some_name","value":%v}`, value)
				reader := bytes.NewReader([]byte(body))
				req, _ := http.NewRequest("POST", "/", reader)

				server.httpHandler.ServeHTTP(w, req)
				assert.Equal(t, 200, w.Code)
			}

			categoricalValues := []string{"a", "b", "a", "de", "c", "b", "a"}
			for _, value := range categoricalValues {
				w := httptest.NewRecorder()
				body := fmt.Sprintf(`{"name":"some_other_name","value":"%v"}`, value)
				reader := bytes.NewReader([]byte(body))
				req, _ := http.NewRequest("POST", "/", reader)

				server.httpHandler.ServeHTTP(w, req)
				assert.Equal(t, 200, w.Code)
			}

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/", nil)
			server.httpHandler.ServeHTTP(w, req)
			assert.Equal(t, 200, w.Code)

			response := map[string][]map[string]interface{}{}
			body, err := ioutil.ReadAll(w.Body)
			require.Nil(t, err)

			err = json.Unmarshal(body, &response)
			require.Nil(t, err)

			numericalResponseValues := []float64{}
			for _, measurement := range response["some_name"] {
				numericalResponseValues = append(numericalResponseValues, measurement["value"].(float64))
			}

			assert.ElementsMatch(t, numericalResponseValues, numericalValues)

			categoricalResponseValues := []string{}
			for _, measurement := range response["some_other_name"] {
				categoricalResponseValues = append(categoricalResponseValues, measurement["value"].(string))
			}

			assert.ElementsMatch(t, categoricalResponseValues, categoricalValues)
		})
	})
}

func Test_ServerParams(t *testing.T) {
	formerDataPath := dataPath
	dataPath = "test_data"
	defer func() {
		os.RemoveAll(dataPath)
		dataPath = formerDataPath
	}()
	server := NewServer(ServerConfig{MemorySize: 24 * 1024 * 1024, DiskSize: 24 * 1024 * 1024})

	t.Run("HTTPHandler", func(t *testing.T) {
		t.Run("POSTing and GETing measurements works", func(t *testing.T) {
			numericalValues := []float64{60, 10, 40, 20, 50, 42}
			for i, value := range numericalValues {
				w := httptest.NewRecorder()
				body := fmt.Sprintf(`{"name":"some_name","value":%v, "timestamp":%v}`, value, i*1000+1000)
				reader := bytes.NewReader([]byte(body))
				req, _ := http.NewRequest("POST", "/", reader)

				server.httpHandler.ServeHTTP(w, req)
				assert.Equal(t, 200, w.Code)
			}

			categoricalValues := []string{"a", "b", "a", "de", "c", "b", "a"}
			for i, value := range categoricalValues {
				w := httptest.NewRecorder()
				body := fmt.Sprintf(`{"name":"some_other_name","value":"%v", "timestamp":%v}`, value, i*1000+1000)
				reader := bytes.NewReader([]byte(body))
				req, _ := http.NewRequest("POST", "/", reader)

				server.httpHandler.ServeHTTP(w, req)
				assert.Equal(t, 200, w.Code)
			}

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/?start=2001", nil)
			server.httpHandler.ServeHTTP(w, req)
			assert.Equal(t, 200, w.Code)

			response := map[string][]map[string]interface{}{}
			body, err := ioutil.ReadAll(w.Body)
			require.Nil(t, err)
			err = json.Unmarshal(body, &response)
			require.Nil(t, err)

			numericalResponseValues := []float64{}
			for _, measurement := range response["some_name"] {
				numericalResponseValues = append(numericalResponseValues, measurement["value"].(float64))
			}

			assert.ElementsMatch(t, numericalResponseValues, numericalValues[2:])

			categoricalResponseValues := []string{}
			for _, measurement := range response["some_other_name"] {
				categoricalResponseValues = append(categoricalResponseValues, measurement["value"].(string))
			}

			assert.ElementsMatch(t, categoricalResponseValues, categoricalValues[2:])
		})
	})
}
