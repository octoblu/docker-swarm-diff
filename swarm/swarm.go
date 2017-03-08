package swarm

import (
	"context"
	"sort"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/octoblu/docker-swarm-diff/server"
	De "github.com/visionmedia/go-debug"
)

var debug = De.Debug("docker-swarm-diff:swarm")

// GetServers returns all servers swarm knows about,
// each containing the service instances that swarm thinks
// should be running on that server
func GetServers() ([]server.Server, error) {
	swarmServerMap, err := getSwarmServerMap()
	if err != nil {
		return nil, err
	}

	services, err := getServiceList()
	if err != nil {
		return nil, err
	}

	errChan := make(chan error, len(services))

	for _, service := range services {
		go func(service swarm.Service) {
			instances, err := getInstancesForService(service)
			if err != nil {
				errChan <- err
			}

			for _, instance := range instances {
				serverID := instance.ServerID()
				swarmServerMap[serverID].AddInstance(instance)
			}
			errChan <- nil
		}(service)
	}

	for i := 0; i < len(services); i++ {
		if getInstanceErr := <-errChan; getInstanceErr != nil {
			return nil, getInstanceErr
		}
	}

	return convertToServers(swarmServerMap), nil
}

func convertToServers(swarmServerMap map[string]*Server) []server.Server {
	debug("convertToServers")
	var servers []server.Server

	for _, swarmServer := range swarmServerMap {
		servers = append(servers, swarmServer)
	}

	sort.Slice(servers, func(i, j int) bool {
		return servers[i].String() < servers[j].String()
	})

	return servers
}

func getSwarmServerMap() (map[string]*Server, error) {
	debug("getSwarmServerMap")
	swarmServerMap := make(map[string]*Server)

	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	nodes, err := cli.NodeList(context.Background(), types.NodeListOptions{})
	if err != nil {
		return nil, err
	}

	for _, node := range nodes {
		swarmServerMap[node.ID] = NewServer(node)
	}

	return swarmServerMap, nil
}

func getInstancesForService(service swarm.Service) ([]*Instance, error) {
	debug("getInstancesForService")
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	filter := filters.NewArgs()
	filter.Add("service", service.ID)

	tasks, err := cli.TaskList(context.Background(), types.TaskListOptions{Filters: filter})
	if err != nil {
		return nil, err
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Meta.CreatedAt.Unix() > tasks[j].Meta.CreatedAt.Unix()
	})

	knownInstances := make(map[string]bool)
	instances := make([]*Instance, 0)
	for _, task := range tasks {
		instance := NewInstance(service.Spec.Name, task)

		if _, ok := knownInstances[instance.Key()]; ok {
			continue
		}

		instances = append(instances, instance)
		knownInstances[instance.Key()] = true
	}

	return instances, nil
}

func getServiceList() ([]swarm.Service, error) {
	debug("getServiceList")
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	return cli.ServiceList(context.Background(), types.ServiceListOptions{})
}
