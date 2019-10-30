package mhist

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/alexmorten/mhist/models"
	"github.com/rs/cors"
)

//HTTPHandler handles http connections
type HTTPHandler struct {
	Server     *Server
	Port       int
	httpServer *http.Server

	corsHandler http.Handler
}

//Init the http mux
func (h *HTTPHandler) Init() {
	mux := http.NewServeMux()
	mux.HandleFunc("/meta", h.serveStoredMeta)
	mux.Handle("/", h)
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
	})
	h.corsHandler = c.Handler(mux)
}

//Run the handler
func (h *HTTPHandler) Run() {
	h.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%v", h.Port),
		Handler: h.corsHandler,
	}

	fmt.Println("http_handler running on ", h.httpServer.Addr)
	err := h.httpServer.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

//Shutdown the HTTPHandler
func (h *HTTPHandler) Shutdown() {
	if h.httpServer != nil {
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		h.httpServer.Shutdown(timeoutCtx)
	}
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()
	switch r.Method {
	case http.MethodPost:
		h.handlePost(w, r)
	case http.MethodGet:
		h.handleGet(w, r)
	}
}

func (h *HTTPHandler) serveStoredMeta(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

	infos := h.Server.store.GetStoredMetaInfo()
	byteSlice, err := json.Marshal(infos)
	if err != nil {
		renderError(err, w, http.StatusInternalServerError)
		return
	}

	w.Write(byteSlice)
}

func (h *HTTPHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	byteSlice, err := ioutil.ReadAll(r.Body)
	if err != nil {
		renderError(err, w, http.StatusBadRequest)
		return
	}
	h.Server.HandleNewMessage(byteSlice, func(err error, status int) {
		renderError(err, w, status)
	})
}

type getParams struct {
	startTs          int64
	endTs            int64
	filterDefinition models.FilterDefinition
}

func (h *HTTPHandler) handleGet(w http.ResponseWriter, r *http.Request) {
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

	responseMap := h.Server.store.GetMeasurementsInTimeRange(params.startTs, params.endTs, params.filterDefinition)
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
	granularityParam := params.Get("granularity")
	namesParam := params.Get("names")
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

	if granularityParam != "" {
		granularity, err := time.ParseDuration(granularityParam)
		if err != nil {
			return nil, err
		}
		p.filterDefinition.Granularity = granularity
	}

	if namesParam != "" {
		names := strings.Split(namesParam, ",")
		p.filterDefinition.Names = names
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
		w.WriteHeader(status)
		w.Write(data)
	}
}
