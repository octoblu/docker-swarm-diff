package helpers

import (
	"github.com/octoblu/docker-swarm-diff/reality"
	"github.com/octoblu/docker-swarm-diff/server"
	"github.com/octoblu/docker-swarm-diff/swarm"
)

// Result is meant to be passed through a channel in
// an Get*Async method
type Result struct {
	Servers []server.Server
	Err     error
}

// GetExpectationsAsync retrieves the expectations in the Background
// and returns a channel that will eventually contain a result
func GetExpectationsAsync() chan Result {
	resultChan := make(chan Result)

	go func() {
		servers, err := swarm.GetServers()
		resultChan <- Result{Servers: servers, Err: err}
	}()

	return resultChan
}

// GetRealityAsync retrieves the reality in the Background
// and returns a channel that will eventually contain a result
func GetRealityAsync() chan Result {
	resultChan := make(chan Result)

	go func() {
		servers, err := reality.GetServers()
		resultChan <- Result{Servers: servers, Err: err}
	}()

	return resultChan
}
