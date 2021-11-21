package server

import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
)

type httpServer struct {
	Log *Log
}

type ProduceRequest struct {
	Record Record
}

type ProduceResponse struct {
	Offset uint64
}

type ConsumeRequest struct {
	Offset uint64
}

type ConsumeResponse struct {
	Record Record
}

func newHttpServer() *httpServer {
	return &httpServer{
		Log: NewLog(),
	}
}

func NewHttpServer(addr string) *http.Server {
	httpsrv := newHttpServer()
	r := mux.NewRouter()
	r.HandleFunc("/", httpsrv.handleProduce).Methods("POST")
	r.HandleFunc("/", httpsrv.handleConsume).Methods("GET")
	return &http.Server{
		Addr: addr,
		Handler: r,
	}
}

func (s *httpServer) handleProduce(w http.ResponseWriter, r *http.Request) {
	var req ProduceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if (err != nil) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	off, err := s.Log.Append(req.Record)
	if (err != nil) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res := ProduceResponse{Offset: off}
	err = json.NewEncoder(w).Encode(res)
	if (err != nil) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *httpServer) handleConsume(w http.ResponseWriter, r *http.Request) {
	var req ConsumeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if (err != nil) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rec, err := s.Log.Read(req.Offset)
	if (err == ErrOffsetNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if (err != nil) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res := ConsumeResponse{Record: rec}
	err = json.NewEncoder(w).Encode(res)
	if (err != nil) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}


