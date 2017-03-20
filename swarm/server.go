package swarm

import (
	"sort"

	"github.com/docker/engine-api/types/swarm"
	"github.com/octoblu/docker-swarm-diff/server"
)

// Server represent what a server should looks like
// according to swarm
type Server struct {
	node      swarm.Node
	instances []*Instance
}

// NewServer constructs a new Server instance from
// a docker node
func NewServer(node swarm.Node) *Server {
	var instances []*Instance

	return &Server{
		node:      node,
		instances: instances,
	}
}

// AddInstance adds an instance to this server
func (serv *Server) AddInstance(instance *Instance) {
	serv.instances = append(serv.instances, instance)
}

// ServiceInstances represent what service instances
// swarm believes to be running on this particular service
func (serv *Server) ServiceInstances() ([]server.ServiceInstance, error) {
	serviceInstances := make([]server.ServiceInstance, len(serv.instances))
	for i, instance := range serv.instances {
		serviceInstances[i] = instance
	}

	sort.Slice(serviceInstances, func(i, j int) bool {
		return serviceInstances[i].String() < serviceInstances[j].String()
	})
	return serviceInstances, nil
}

// String returns a string representation of the server
func (serv *Server) String() string {
	return serv.node.Description.Hostname
}
