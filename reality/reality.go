package reality

import (
	"context"
	"sort"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/octoblu/docker-swarm-diff/server"
)

// GetServers returns all servers swarm knows about,
// each containing the service instances that are actually
// running on that server (according to that server)
func GetServers() ([]server.Server, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	nodes, err := cli.NodeList(context.Background(), types.NodeListOptions{})
	if err != nil {
		return nil, err
	}

	servers := make([]server.Server, len(nodes))
	for i, node := range nodes {
		servers[i] = NewServer(node)
	}

	sort.Slice(servers, func(i, j int) bool {
		return servers[i].String() < servers[j].String()
	})

	return servers, nil
}
