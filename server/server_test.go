package server

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/albertogviana/docker-swarm-service-status/service"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ServerTestSuite struct {
	suite.Suite
}

const DockerHost = "unix:///var/run/docker.sock"
const DockerAPIVersion = "v1.33"

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

func (s *ServerTestSuite) Test_NewServer_ReturnServer() {
	serviceMock := new(ServiceMock)
	server := NewServer(serviceMock)

	assert.IsType(s.T(), &Server{}, server)
}

func (s *ServerTestSuite) Test_DeploymentStatus_ReturnSuccess() {
	serviceMock := new(ServiceMock)

	taskStatus := []service.TaskStatus{}

	ts := service.TaskStatus{
		"evv1jw9o7981mrp0p50j1gy5k",
		time.Date(2017, time.November, 26, 21, 47, 35, 0, time.UTC),
		"running",
		"running",
		"started",
		"",
		"albertogviana/docker-routing-mesh:1.0.0@sha256:87e5c74f8042848893440b24a33ea0e3494b9da475987b0e704f0d3262bce3cd",
	}

	taskStatus = append(taskStatus, ts)

	replicas := uint64(1)

	deploymentStatusMock := service.ServiceStatus{
		"tt3otdsnkd1kgh80u45bwmcb4",
		"docker-routing-mesh",
		"",
		taskStatus,
		&replicas,
		1,
		0,
		nil,
	}

	data, _ := json.Marshal(deploymentStatusMock)

	serviceName := "docker-routing-mesh"
	image := "albertogviana/docker-routing-mesh:1.0.0"

	serviceMock.On("GetDeploymentStatus", serviceName, image).Return(deploymentStatusMock, nil)
	server := &Server{
		serviceMock,
	}

	muxRouter := mux.NewRouter().StrictSlash(true).UseEncodedPath()
	router(muxRouter, server)

	imageByte := base64.URLEncoding.EncodeToString([]byte(image))

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/docker-swarm-service-status/deployment-status/%s/%s", serviceName, imageByte), nil)

	muxRouter.ServeHTTP(rec, req)

	s.Equal(200, rec.Code)

	s.Equal(string(data), rec.Body.String())
}

func (s *ServerTestSuite) Test_DeploymentStatus_InvalidBase64Parameter() {
	serviceMock := new(ServiceMock)

	serviceName := "docker-routing-mesh"
	image := "docker-routing-mesh:1.0.0"

	server := &Server{
		serviceMock,
	}

	muxRouter := mux.NewRouter().StrictSlash(true).UseEncodedPath()
	router(muxRouter, server)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/docker-swarm-service-status/deployment-status/%s/%s", serviceName, image), nil)

	muxRouter.ServeHTTP(rec, req)

	s.Equal(400, rec.Code)

	s.Equal("{\"error\": \"Invalid base64 encode for the image parameter.\"}", rec.Body.String())
}

func (s *ServerTestSuite) Test_DeploymentStatus_ReturnError() {
	serviceMock := new(ServiceMock)

	serviceName := "docker-routing-mesh"
	image := "albertogviana/docker-routing-mesh:1.0.0"

	serviceMock.On("GetDeploymentStatus", serviceName, image).Return(service.ServiceStatus{}, errors.New("Not able to connect on unix:///var/run/docker.sock"))
	server := &Server{
		serviceMock,
	}

	muxRouter := mux.NewRouter().StrictSlash(true).UseEncodedPath()
	router(muxRouter, server)

	imageByte := base64.URLEncoding.EncodeToString([]byte(image))

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/docker-swarm-service-status/deployment-status/%s/%s", serviceName, imageByte), nil)

	muxRouter.ServeHTTP(rec, req)

	s.Equal(500, rec.Code)

	s.Equal("{\"error\": \"Not able to connect on unix:///var/run/docker.sock\"}", rec.Body.String())
}

func (s *ServerTestSuite) Test_ServiceStatus_ReturnError() {
	serviceMock := new(ServiceMock)

	serviceName := "docker-routing-mesh"

	serviceMock.On("GetServiceStatus", serviceName).Return(service.ServiceStatus{}, errors.New("Not able to connect on unix:///var/run/docker.sock"))
	server := &Server{
		serviceMock,
	}

	muxRouter := mux.NewRouter().StrictSlash(true).UseEncodedPath()
	router(muxRouter, server)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/docker-swarm-service-status/service-status/%s", serviceName), nil)

	muxRouter.ServeHTTP(rec, req)

	s.Equal(500, rec.Code)

	s.Equal("{\"error\": \"Not able to connect on unix:///var/run/docker.sock\"}", rec.Body.String())
}

func (s *ServerTestSuite) Test_ServiceStatus_ReturnSuccess() {
	serviceMock := new(ServiceMock)

	taskStatus := []service.TaskStatus{}

	ts := service.TaskStatus{
		"evv1jw9o7981mrp0p50j1gy5k",
		time.Date(2017, time.November, 26, 21, 47, 35, 0, time.UTC),
		"running",
		"running",
		"started",
		"",
		"albertogviana/docker-routing-mesh:1.0.0@sha256:87e5c74f8042848893440b24a33ea0e3494b9da475987b0e704f0d3262bce3cd",
	}

	taskStatus = append(taskStatus, ts)

	replicas := uint64(1)

	deploymentStatusMock := service.ServiceStatus{
		"tt3otdsnkd1kgh80u45bwmcb4",
		"docker-routing-mesh",
		"",
		taskStatus,
		&replicas,
		1,
		0,
		nil,
	}

	data, _ := json.Marshal(deploymentStatusMock)

	serviceName := "docker-routing-mesh"

	serviceMock.On("GetServiceStatus", serviceName).Return(deploymentStatusMock, nil)
	server := &Server{
		serviceMock,
	}

	muxRouter := mux.NewRouter().StrictSlash(true).UseEncodedPath()
	router(muxRouter, server)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/docker-swarm-service-status/service-status/%s", serviceName), nil)

	muxRouter.ServeHTTP(rec, req)

	s.Equal(200, rec.Code)

	s.Equal(string(data), rec.Body.String())
}

type ServiceMock struct {
	mock.Mock
}

func (s *ServiceMock) GetDeploymentStatus(serviceName string, image string) (service.ServiceStatus, error) {
	args := s.Called(serviceName, image)
	return args.Get(0).(service.ServiceStatus), args.Error(1)
}

func (s *ServiceMock) GetServiceStatus(serviceName string) (service.ServiceStatus, error) {
	args := s.Called(serviceName)
	return args.Get(0).(service.ServiceStatus), args.Error(1)
}

func (s *ServiceMock) GetService(filter filters.Args) (swarm.Service, error) {
	args := s.Called(filter)
	return args.Get(0).(swarm.Service), args.Error(1)
}

func (s *ServiceMock) GetTask(filter filters.Args) ([]swarm.Task, error) {
	args := s.Called(filter)
	return args.Get(0).([]swarm.Task), args.Error(1)
}
