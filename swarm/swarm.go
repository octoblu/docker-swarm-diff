package swarm

import "github.com/octoblu/docker-swarm-diff/server"

// GetServers returns all servers swarm knows about,
// each containing the service instances that swarm thinks
// should be running on that server
func GetServers() ([]server.Server, error) {
	return nil, nil
}
