package reality

import (
	"context"
	"sort"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/octoblu/docker-swarm-diff/server"
	De "github.com/visionmedia/go-debug"
)

var debug = De.Debug("docker-swarm-diff:reality")

// GetServers returns all servers swarm knows about,
// each containing the service instances that are actually
// running on that server (according to that server)
func GetServers() ([]server.Server, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	debug("cli.NodeList")
	nodes, err := cli.NodeList(context.Background(), types.NodeListOptions{})
	if err != nil {
		return nil, err
	}

	errChan := make(chan error, len(nodes))
	servers := make([]server.Server, len(nodes))

	for i, node := range nodes {
		serv := NewServer(node)
		servers[i] = serv
		go func() {
			debug("serv.FetchServerInstaces: %v", serv.String())
			errChan <- serv.FetchServerInstances()
			debug("serv.FetchServerInstaces (done): %v, %v", serv.String())
		}()
	}

	debug("range errChan")
	for i := 0; i < len(servers); i++ {
		if fetchErr := <-errChan; fetchErr != nil {
			return nil, fetchErr
		}
	}

	debug("sort.Slice")
	sort.Slice(servers, func(i, j int) bool {
		return servers[i].String() < servers[j].String()
	})

	return servers, nil
}
