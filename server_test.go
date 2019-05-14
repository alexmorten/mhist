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

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Server(t *testing.T) {
	formerDataPath := dataPath
	dataPath = "test_data"
	defer func() {
		os.RemoveAll(dataPath)
		dataPath = formerDataPath
	}()
	server := NewServer(ServerConfig{MemorySize: 24 * 1024 * 1024, DiskSize: 24 * 1024 * 1024})

	Convey("HTTPHandler", t, func() {
		Convey("POSTing and GETing measurements works", func() {
			numericalValues := []float64{10, 60, 40, 20, 50, 42}
			for _, value := range numericalValues {
				w := httptest.NewRecorder()
				body := fmt.Sprintf(`{"name":"some_name","value":%v}`, value)
				reader := bytes.NewReader([]byte(body))
				req, _ := http.NewRequest("POST", "/", reader)

				server.httpHandler.ServeHTTP(w, req)
				So(w.Code, ShouldEqual, 200)
			}

			categoricalValues := []string{"a", "b", "a", "de", "c", "b", "a"}
			for _, value := range categoricalValues {
				w := httptest.NewRecorder()
				body := fmt.Sprintf(`{"name":"some_other_name","value":"%v"}`, value)
				reader := bytes.NewReader([]byte(body))
				req, _ := http.NewRequest("POST", "/", reader)

				server.httpHandler.ServeHTTP(w, req)
				So(w.Code, ShouldEqual, 200)
			}

			// commit diskStore to avoid sleeping
			server.store.diskStore.commit()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/", nil)
			server.httpHandler.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, 200)

			response := map[string][]map[string]interface{}{}
			body, err := ioutil.ReadAll(w.Body)
			So(err, ShouldBeNil)
			err = json.Unmarshal(body, &response)
			So(err, ShouldBeNil)
			So(len(response["some_name"]), ShouldEqual, len(numericalValues))
			So(len(response["some_other_name"]), ShouldEqual, len(categoricalValues))

			numericalResponseValues := []float64{}
			for _, measurement := range response["some_name"] {
				numericalResponseValues = append(numericalResponseValues, measurement["value"].(float64))
			}

			So(numericalResponseValues, ShouldResemble, numericalValues)

			categoricalResponseValues := []string{}
			for _, measurement := range response["some_other_name"] {
				categoricalResponseValues = append(categoricalResponseValues, measurement["value"].(string))
			}

			So(categoricalResponseValues, ShouldResemble, categoricalValues)
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

	Convey("HTTPHandler", t, func() {
		Convey("POSTing and GETing measurements works", func() {
			numericalValues := []float64{60, 10, 40, 20, 50, 42}
			for i, value := range numericalValues {
				w := httptest.NewRecorder()
				body := fmt.Sprintf(`{"name":"some_name","value":%v, "timestamp":%v}`, value, i*1000+1000)
				reader := bytes.NewReader([]byte(body))
				req, _ := http.NewRequest("POST", "/", reader)

				server.httpHandler.ServeHTTP(w, req)
				So(w.Code, ShouldEqual, 200)
			}

			categoricalValues := []string{"a", "b", "a", "de", "c", "b", "a"}
			for i, value := range categoricalValues {
				w := httptest.NewRecorder()
				body := fmt.Sprintf(`{"name":"some_other_name","value":"%v", "timestamp":%v}`, value, i*1000)
				reader := bytes.NewReader([]byte(body))
				req, _ := http.NewRequest("POST", "/", reader)

				server.httpHandler.ServeHTTP(w, req)
				So(w.Code, ShouldEqual, 200)
			}

			// commit diskStore to avoid sleeping
			server.store.diskStore.commit()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/?start=2001", nil)
			server.httpHandler.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, 200)

			response := map[string][]map[string]interface{}{}
			body, err := ioutil.ReadAll(w.Body)
			fmt.Println(string(body))
			So(err, ShouldBeNil)
			err = json.Unmarshal(body, &response)
			So(err, ShouldBeNil)
			So(len(response["some_name"]), ShouldEqual, len(numericalValues)-2)
			So(len(response["some_other_name"]), ShouldEqual, len(categoricalValues)-2)

			numericalResponseValues := []float64{}
			for _, measurement := range response["some_name"] {
				numericalResponseValues = append(numericalResponseValues, measurement["value"].(float64))
			}

			So(numericalResponseValues, ShouldResemble, numericalValues[2:])

			categoricalResponseValues := []string{}
			for _, measurement := range response["some_other_name"] {
				categoricalResponseValues = append(categoricalResponseValues, measurement["value"].(string))
			}

			So(categoricalResponseValues, ShouldResemble, categoricalValues[2:])
		})

	})
}
