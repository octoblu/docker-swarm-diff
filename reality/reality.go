package reality

import "github.com/octoblu/docker-swarm-diff/server"

// GetServers returns all servers swarm knows about,
// each containing the service instances that are actually
// running on that server (according to that server)
func GetServers() ([]server.Server, error) {
	return nil, nil
}
