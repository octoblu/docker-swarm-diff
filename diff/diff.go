package diff

import "github.com/octoblu/docker-swarm-diff/server"

// Diff defines a difference between expectation and reality
type Diff interface {
	// String returns a formatted difference
	String() string
}

// Differentiate returns an array of differences between
// expecations and reality
func Differentiate(expecations, reality []server.Server) []Diff {
	return nil
}
