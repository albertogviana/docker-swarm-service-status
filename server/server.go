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

// NewServer returns a new instance of the Server structure
func NewServer(service service.Services) *Server {
	return &Server{
		service,
	}
}

// Run bootstrap the server
func (s *Server) Run() {
	log.Println("Docker Deployment Status Starting")

	r := mux.NewRouter().StrictSlash(true).UseEncodedPath()
	deploymentStatusHandler(r, s)

	log.Println("Docker Deployment Status Started")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", r))
}

func deploymentStatusHandler(r *mux.Router, s *Server) {
	r.HandleFunc("/v1/docker-swarm-deployment-status/{service}/{image}", s.DeploymentStatus).Methods("GET")
}

// DeploymentStatus returns the current state of the service
func (s *Server) DeploymentStatus(w http.ResponseWriter, r *http.Request) {
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
