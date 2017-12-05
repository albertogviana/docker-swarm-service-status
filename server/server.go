package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/albertogviana/docker-swarm-deployment-status/service"
	"github.com/gorilla/mux"
)

// Server defined structure
type Server struct {
	Service service.Services
}

//Response message
type Response struct {
	Status string
}

// NewServer returns a new instance of the Server structure
func NewServer(service service.Services) *Server {
	return &Server{
		service,
	}
}

// Run bootstrap the server
func (s *Server) Run() {
	log.Println("Docker Service Status Starting")

	r := mux.NewRouter().StrictSlash(true).UseEncodedPath()
	router(r, s)

	log.Println("Docker Service Status Started")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", r))
}

func router(r *mux.Router, s *Server) {
	r.HandleFunc("/v1/docker-swarm-service-status/service-status/{service}", s.ServiceStatusHandler).Methods("GET")
	r.HandleFunc("/v1/docker-swarm-service-status/deployment-status/{service}/{image}", s.DeploymentStatusHandler).Methods("GET")
	r.HandleFunc("/v1/docker-swarm-service-status/health", s.HealthHandler).Methods("GET")
}

// DeploymentStatusHandler returns the current state of the service
func (s *Server) DeploymentStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	serviceName := vars["service"]
	image := vars["image"]

	imageByte, err := base64.URLEncoding.DecodeString(image)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Invalid base64 encode for the image parameter."}`)
		return
	}

	status, err := s.Service.GetDeploymentStatus(serviceName, string(imageByte))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	js, _ := json.Marshal(status)
	w.Write(js)
}

// ServiceStatusHandler returns the current state of the service
func (s *Server) ServiceStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	serviceName := vars["service"]

	status, err := s.Service.GetServiceStatus(serviceName)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	js, _ := json.Marshal(status)
	w.Write(js)
}

// HealthHandler is used for health checks
func (s *Server) HealthHandler(w http.ResponseWriter, req *http.Request) {
	js, _ := json.Marshal(Response{Status: "OK"})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}
