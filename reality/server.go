package reality

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/docker/machine/commands/mcndirs"
	"github.com/docker/machine/libmachine"
	"github.com/docker/machine/libmachine/check"
	"github.com/octoblu/docker-swarm-diff/server"
)

// Server represents a server and can list the docker containers
// running on it
type Server struct {
	node swarm.Node
}

// NewServer constructs a new server
func NewServer(node swarm.Node) *Server {
	return &Server{node: node}
}

// ServiceInstances represent what service instances are running
// on this server, based off of the docker containers that are running
func (serv *Server) ServiceInstances() ([]server.ServiceInstance, error) {
	containers, err := serv.getContainers()
	if err != nil {
		return nil, err
	}
	instances := make([]server.ServiceInstance, len(containers))
	for i, container := range containers {
		instances[i] = NewInstance(container)
	}

	sort.Slice(instances, func(i, j int) bool {
		return instances[i].String() < instances[j].String()
	})
	return instances, nil
}

// String returns a string representation of the server
func (serv *Server) String() string {
	return serv.node.Description.Hostname
}

func (serv *Server) getContainers() ([]types.Container, error) {
	var returnErr error
	var containers []types.Container

	inEnvSandboxDo(func() {
		dockerMachineEnv, err := serv.getDockerMachineEnv()
		if err != nil {
			returnErr = err
			return
		}

		setEnv(dockerMachineEnv)
		cli, err := client.NewEnvClient()
		if err != nil {
			returnErr = err
			return
		}

		containers, returnErr = cli.ContainerList(context.Background(), types.ContainerListOptions{})
	})
	return containers, returnErr
}

func (serv *Server) getDockerMachineEnv() (map[string]string, error) {
	cli := libmachine.NewClient(mcndirs.GetBaseDir(), os.Getenv("DOCKER_CERT_PATH"))
	host, err := cli.Load(serv.node.Description.Hostname)
	if err != nil {
		return nil, err
	}

	dockerHost, _, err := check.DefaultConnChecker.Check(host, false)
	if err != nil {
		return nil, err
	}

	dockerMachineEnv := make(map[string]string)
	dockerMachineEnv["DOCKER_CERT_PATH"] = filepath.Join(mcndirs.GetMachineDir(), host.Name)
	dockerMachineEnv["DOCKER_HOST"] = dockerHost
	dockerMachineEnv["DOCKER_TLS_VERIFY"] = "1"
	dockerMachineEnv["DOCKER_MACHINE_NAME"] = host.Name
	return dockerMachineEnv, nil
}

func inEnvSandboxDo(fn func()) {
	envBackup := make(map[string]string)

	for _, pair := range os.Environ() {
		key := strings.Split(pair, "=")[0]
		envBackup[key] = os.Getenv(key)
	}

	fn()

	for _, pair := range os.Environ() {
		key := strings.Split(pair, "=")[0]
		os.Unsetenv(key)
	}

	setEnv(envBackup)
}

func setEnv(envMap map[string]string) {
	for key, value := range envMap {
		os.Setenv(key, value)
	}
}
