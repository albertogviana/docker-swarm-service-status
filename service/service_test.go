package service

import (
	"os/exec"
	"testing"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
}

const DockerHost = "unix:///var/run/docker.sock"
const DockerAPIVersion = "v1.33"

func TestServiceTestSuite(t *testing.T) {
	createTestServices()
	suite.Run(t, new(ServiceTestSuite))
	removeTestServices()
}

func (s *ServiceTestSuite) Test_NewService_ReturnService() {
	service := NewService(DockerHost, DockerAPIVersion, map[string]string{})

	assert.NotNil(s.T(), service)
	assert.Equal(s.T(), DockerHost, service.Host)
}

func (s *ServiceTestSuite) Test_GetService_ReturnSwarmService() {
	serviceName := "docker-routing-mesh"
	filterList := filters.NewArgs()
	filterList.Add("name", serviceName)

	service := NewService(DockerHost, DockerAPIVersion, map[string]string{})
	swarmService, err := service.GetService(filterList)

	assert.NoError(s.T(), err)
	assert.IsType(s.T(), swarm.Service{}, swarmService)
	assert.Equal(s.T(), serviceName, swarmService.Spec.Name)
}

func (s *ServiceTestSuite) Test_GetService_ReturnEmptySwarmService_NoServiceFound() {
	serviceName := "my-service"
	filterList := filters.NewArgs()
	filterList.Add("name", serviceName)

	service := NewService(DockerHost, DockerAPIVersion, map[string]string{})
	swarmService, err := service.GetService(filterList)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), swarm.Service{}, swarmService)
}

func (s *ServiceTestSuite) Test_GetService_ReturnError_InvalidFilter() {
	serviceName := "docker-routing-mesh"
	filterList := filters.NewArgs()
	filterList.Add("invalidFilter", serviceName)

	service := NewService(DockerHost, DockerAPIVersion, map[string]string{})
	_, err := service.GetService(filterList)

	assert.Error(s.T(), err, "Error response from daemon: {\"message\":\"Invalid filter 'invalidFilter'\"}")
}

func (s *ServiceTestSuite) Test_GetTask_ReturnSwarmTask() {
	serviceName := "docker-routing-mesh"
	filterList := filters.NewArgs()
	filterList.Add("name", serviceName)

	service := NewService(DockerHost, DockerAPIVersion, map[string]string{})
	swarmService, err := service.GetService(filterList)

	assert.NoError(s.T(), err)

	tasksFilter := filters.NewArgs()
	tasksFilter.Add("service", swarmService.ID)

	swarmTask, err := service.GetTask(tasksFilter)

	assert.NoError(s.T(), err)
	assert.IsType(s.T(), []swarm.Task{}, swarmTask)
}

func (s *ServiceTestSuite) Test_GetTask_ReturnError_InvalidFilter() {
	serviceName := "docker-routing-mesh"
	filterList := filters.NewArgs()
	filterList.Add("name", serviceName)

	service := NewService(DockerHost, DockerAPIVersion, map[string]string{})
	swarmService, err := service.GetService(filterList)

	assert.NoError(s.T(), err)

	tasksFilter := filters.NewArgs()
	tasksFilter.Add("invalidFilter", swarmService.ID)

	_, err = service.GetTask(tasksFilter)

	assert.Error(s.T(), err, "Error response from daemon: {\"message\":\"Invalid filter 'invalidFilter'\"}")
}

func (s *ServiceTestSuite) Test_GetDeploymentStatus_ReturnServiceStatus_ServiceNotExists() {
	service := NewService(DockerHost, DockerAPIVersion, map[string]string{})
	deploymentStatus, err := service.GetDeploymentStatus("my-service", "my-image:1.0.0")

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "The my-service service was not found in the cluster.", deploymentStatus.Err)
}

func (s *ServiceTestSuite) Test_GetDeploymentStatus_ReturnDeploymentStatus_RunningDifferentImage() {
	service := NewService(DockerHost, DockerAPIVersion, map[string]string{})
	serviceName := "docker-routing-mesh"
	deploymentStatus, err := service.GetDeploymentStatus(serviceName, "albertogviana/docker-routing-mesh:1.0.1")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), deploymentStatus.ID)
	assert.Equal(s.T(), serviceName, deploymentStatus.Name)
	assert.Equal(s.T(), "The albertogviana/docker-routing-mesh:1.0.1 image was not deployed or not found in the current tasks running.", deploymentStatus.Err)
}

func (s *ServiceTestSuite) Test_GetDeploymentStatus_ReturnServiceStatus() {
	defer func() {
		exec.Command("docker", "service", "scale", "docker-routing-mesh=1").Output()
	}()

	service := NewService(DockerHost, DockerAPIVersion, map[string]string{})
	serviceName := "docker-routing-mesh"
	image := "albertogviana/docker-routing-mesh:1.0.0"
	deploymentStatus, err := service.GetDeploymentStatus(serviceName, image)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), deploymentStatus.ID)
	assert.Equal(s.T(), serviceName, deploymentStatus.Name)
	assert.Equal(s.T(), uint64(1), *deploymentStatus.Replicas)
	assert.Equal(s.T(), 1, deploymentStatus.RunningReplicas)
	assert.Equal(s.T(), 0, deploymentStatus.FailedReplicas)
	assert.Nil(s.T(), deploymentStatus.UpdateStatus)

	exec.Command("docker", "service", "scale", "docker-routing-mesh=2").Output()
	deploymentStatus2, err := service.GetDeploymentStatus(serviceName, image)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), deploymentStatus2.ID)
	assert.Equal(s.T(), serviceName, deploymentStatus2.Name)
	assert.Equal(s.T(), uint64(2), *deploymentStatus2.Replicas)
	assert.Equal(s.T(), 2, deploymentStatus2.RunningReplicas)
	assert.Equal(s.T(), 0, deploymentStatus2.FailedReplicas)
	assert.Nil(s.T(), deploymentStatus2.UpdateStatus)
}

func (s *ServiceTestSuite) Test_GetDeploymentStatus_ReturnServiceStatusWithUpdateStatus() {
	defer func() {
		removeTestServices()
		createTestServices()
	}()

	exec.Command("docker", "service", "update", "--image", "albertogviana/docker-routing-mesh:2.0.0", "--replicas", "2", "docker-routing-mesh").Output()

	service := NewService(DockerHost, DockerAPIVersion, map[string]string{})
	serviceName := "docker-routing-mesh"
	deploymentStatus, err := service.GetDeploymentStatus(serviceName, "albertogviana/docker-routing-mesh:2.0.0")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), deploymentStatus.ID)
	assert.Equal(s.T(), serviceName, deploymentStatus.Name)
	assert.Equal(s.T(), uint64(2), *deploymentStatus.Replicas)
	assert.Equal(s.T(), 2, deploymentStatus.RunningReplicas)
	assert.Equal(s.T(), 0, deploymentStatus.FailedReplicas)
	assert.NotNil(s.T(), deploymentStatus.UpdateStatus)
	assert.Equal(s.T(), swarm.UpdateStateCompleted, deploymentStatus.UpdateStatus.State)
	assert.Equal(s.T(), "update completed", deploymentStatus.UpdateStatus.Message)
}

func (s *ServiceTestSuite) Test_GetDeploymentStatus_ReturnServiceStatusWithUpdateStatusFailed() {
	defer func() {
		removeTestServices()
		createTestServices()
	}()

	exec.Command("docker", "service", "update", "--image", "albertogviana/docker-routing-mesh:error", "--replicas", "2", "docker-routing-mesh").Output()

	service := NewService(DockerHost, DockerAPIVersion, map[string]string{})
	serviceName := "docker-routing-mesh"
	deploymentStatus, err := service.GetDeploymentStatus(serviceName, "albertogviana/docker-routing-mesh:error")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), deploymentStatus.ID)
	assert.Equal(s.T(), serviceName, deploymentStatus.Name)
	assert.Equal(s.T(), uint64(2), *deploymentStatus.Replicas)
	assert.Equal(s.T(), 0, deploymentStatus.RunningReplicas)

	failedReplicas := false
	if deploymentStatus.FailedReplicas > 0 {
		failedReplicas = true
	}

	assert.True(s.T(), failedReplicas)
	assert.NotNil(s.T(), deploymentStatus.UpdateStatus)
	assert.Equal(s.T(), swarm.UpdateStatePaused, deploymentStatus.UpdateStatus.State)
}

func (s *ServiceTestSuite) Test_GetGetServiceStatus_ReturnServiceNotExists() {
	service := NewService(DockerHost, DockerAPIVersion, map[string]string{})
	deploymentStatus, err := service.GetServiceStatus("my-service")

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "The my-service service was not found in the cluster.", deploymentStatus.Err)
}

func (s *ServiceTestSuite) Test_GetServiceStatus_ReturnServiceStatus() {
	defer func() {
		exec.Command("docker", "service", "scale", "docker-routing-mesh=1").Output()
	}()

	service := NewService(DockerHost, DockerAPIVersion, map[string]string{})
	serviceName := "docker-routing-mesh"
	deploymentStatus, err := service.GetServiceStatus(serviceName)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), deploymentStatus.ID)
	assert.Equal(s.T(), serviceName, deploymentStatus.Name)
	assert.Equal(s.T(), uint64(1), *deploymentStatus.Replicas)
	assert.Equal(s.T(), 1, deploymentStatus.RunningReplicas)
	assert.Nil(s.T(), deploymentStatus.UpdateStatus)

	exec.Command("docker", "service", "scale", "docker-routing-mesh=2").Output()
	deploymentStatus2, err := service.GetServiceStatus(serviceName)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), deploymentStatus2.ID)
	assert.Equal(s.T(), serviceName, deploymentStatus2.Name)
	assert.Equal(s.T(), uint64(2), *deploymentStatus2.Replicas)
	assert.Equal(s.T(), 2, deploymentStatus2.RunningReplicas)
	assert.Nil(s.T(), deploymentStatus2.UpdateStatus)
}

// Util

func createTestServices() {
	createTestService("docker-routing-mesh", []string{}, "", "albertogviana/docker-routing-mesh:1.0.0")
}

func createTestService(name string, labels []string, mode string, image string) {
	args := []string{"service", "create", "--name", name}
	for _, v := range labels {
		args = append(args, "-l", v)
	}
	if len(mode) > 0 {
		args = append(args, "--mode", "global")
	}
	args = append(args, image)
	exec.Command("docker", args...).Output()
}

func removeTestServices() {
	removeTestService("docker-routing-mesh")
}

func removeTestService(name string) {
	exec.Command("docker", "service", "rm", name).Output()
}
