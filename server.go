package mhist

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

//Server is the handler for requests
type Server struct {
	store *Store
	pools *Pools
}

//NewServer returns a new Server
func NewServer(memorySize int) *Server {
	memStore := NewStore(memorySize)
	pools := NewPools(memStore)
	diskStore, err := NewDiskStore(pools)
	if err != nil {
		panic(err)
	}
	memStore.AddSubscriber(diskStore)

	return &Server{
		store: memStore,
		pools: pools,
	}
}

//Run the server
func (s *Server) Run() {
	http.Handle("/", s)
	err := http.ListenAndServe(":6666", nil)
	if err != nil {
		fmt.Println(err)
	}
	s.Shutdown()
}

//Shutdown all goroutines
func (s *Server) Shutdown() {
	s.store.Shutdown()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()
	switch r.Method {
	case http.MethodPost:
		s.handlePost(w, r)
	case http.MethodGet:
		s.handleGet(w, r)
	}
}

type postRequest struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

func (s *Server) handlePost(w http.ResponseWriter, r *http.Request) {
	byteSlice, err := ioutil.ReadAll(r.Body)
	if err != nil {
		renderError(err, w, http.StatusBadRequest)
		return
	}
	data := &postRequest{}
	err = json.Unmarshal(byteSlice, data)
	if err != nil {
		renderError(err, w, http.StatusBadRequest)
		return
	}
	if data.Name == "" {
		err = errors.New("name can't be empty")
		renderError(err, w, http.StatusBadRequest)
		return
	}
	measurement, err := s.constructMeasurementFromPostRequest(data)
	if err != nil {
		renderError(err, w, http.StatusBadRequest)
		return
	}
	s.store.Add(data.Name, measurement)
}

type getParams struct {
	startTs int64
	endTs   int64
}

func (s *Server) handleGet(w http.ResponseWriter, r *http.Request) {
	params, err := parseParams(r.URL.Query())
	if err != nil {
		renderError(err, w, http.StatusBadRequest)
		return
	}
	if params.startTs > params.endTs {
		err := errors.New("start can't be bigger than end")
		renderError(err, w, http.StatusBadRequest)
		return
	}
	responseMap := s.store.GetAllMeasurementsInTimeRange(params.startTs, params.endTs)
	data, err := json.Marshal(responseMap)
	if err != nil {
		renderError(err, w, http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func parseParams(params url.Values) (p *getParams, err error) {
	p = &getParams{}
	startTsParam := params.Get("start")
	endTsParam := params.Get("end")
	if endTsParam == "" {
		p.endTs = time.Now().UnixNano()
	} else {
		p.endTs, err = strconv.ParseInt(endTsParam, 10, 64)
		if err != nil {
			return
		}
	}
	if startTsParam == "" {
		p.startTs = p.endTs - (1 * time.Hour).Nanoseconds()
	} else {
		p.startTs, err = strconv.ParseInt(startTsParam, 10, 64)
	}
	return
}

type errorResponse struct {
	Error string `json:"error"`
}

func renderError(err error, w http.ResponseWriter, status int) {
	fmt.Println(err)
	resp := &errorResponse{Error: err.Error()}
	data, err := json.Marshal(resp)
	if err == nil {
		w.Write(data)
		w.WriteHeader(status)
	}
}

func (s *Server) constructMeasurementFromPostRequest(r *postRequest) (measurement Measurement, err error) {
	switch r.Value.(type) {
	case float64:
		m := s.pools.GetNumericalMeasurement()
		m.Reset()
		m.Ts = time.Now().UnixNano()
		m.Value = r.Value.(float64)
		measurement = m
	case string:
		m := s.pools.GetCategoricalMeasurement()
		m.Reset()
		m.Ts = time.Now().UnixNano()
		m.Value = r.Value.(string)
		measurement = m
	default:
		return nil, errors.New("value is neither a float nor a string")
	}
	return
}
