package main

import (
	"os"

	"github.com/albertogviana/docker-swarm-deployment-status/server"
	"github.com/albertogviana/docker-swarm-deployment-status/service"
)

func main() {

	dockerHost := "unix:///var/run/docker.sock"
	if os.Getenv("DOCKER_HOST") != "" {
		dockerHost = os.Getenv("DOCKER_HOST")
	}

	dockerAPIVersion := "v1.33"
	defaultHeaders := map[string]string{"User-Agent": "docker-swarm-deployment-status-cli-1.0"}

	service := service.NewService(dockerHost, dockerAPIVersion, defaultHeaders)
	server := server.NewServer(service)
	server.Run()
}
