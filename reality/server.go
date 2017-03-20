package reality

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/swarm"
	"github.com/docker/machine/commands/mcndirs"
	"github.com/docker/machine/libmachine"
	"github.com/docker/machine/libmachine/check"
	"github.com/octoblu/docker-swarm-diff/server"
)

// Server represents a server and can list the docker containers
// running on it
type Server struct {
	instances []server.ServiceInstance
	node      swarm.Node
}

// NewServer constructs a new server
func NewServer(node swarm.Node) *Server {
	return &Server{node: node, instances: nil}
}

// FetchServerInstances retrieves the instances off of the server and sets it on
// the server
func (serv *Server) FetchServerInstances() error {
	containers, err := serv.getContainers()
	if err != nil {
		return err
	}

	instances := make([]server.ServiceInstance, len(containers))
	for i, container := range containers {
		instances[i] = NewInstance(container)
	}

	sort.Slice(instances, func(i, j int) bool {
		return instances[i].String() < instances[j].String()
	})

	serv.instances = instances
	return nil
}

// ServiceInstances represent what service instances are running
// on this server, based off of the docker containers that are running
// FetchServerInstances must be called first, else this will return an error
func (serv *Server) ServiceInstances() ([]server.ServiceInstance, error) {
	if serv.instances == nil {
		return nil, fmt.Errorf("FetchServerInstances must be called first (%v)", serv.String())
	}

	return serv.instances, nil
}

// String returns a string representation of the server
func (serv *Server) String() string {
	return serv.node.Description.Hostname
}

func (serv *Server) getContainers() ([]types.Container, error) {
	cli, err := serv.getDockerMachineClient()
	if err != nil {
		return nil, err
	}

	return cli.ContainerList(context.Background(), types.ContainerListOptions{})
}

func (serv *Server) getDockerMachineClient() (client.APIClient, error) {
	cli := libmachine.NewClient(mcndirs.GetBaseDir(), os.Getenv("DOCKER_CERT_PATH"))
	host, err := cli.Load(serv.node.Description.Hostname)
	if err != nil {
		return nil, err
	}

	dockerHost, _, err := check.DefaultConnChecker.Check(host, false)
	if err != nil {
		return nil, err
	}

	dockerCertPath := filepath.Join(mcndirs.GetMachineDir(), host.Name)

	return NewClient(dockerHost, dockerCertPath, true)
}
