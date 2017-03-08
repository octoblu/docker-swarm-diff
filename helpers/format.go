package helpers

import (
	"fmt"

	"github.com/octoblu/docker-swarm-diff/server"
)

// Format pretty prints a slice of servers in the following format:
//
// server-name
// ---------------
// running: docker-proxy.0
// running: traefik.1
//
func Format(servers []server.Server) (string, error) {
	output := ""

	for _, serv := range servers {
		output += formatSubHeading(serv.String())

		instances, err := serv.ServiceInstances()
		if err != nil {
			return "", err
		}

		for _, instance := range instances {
			output += fmt.Sprintln(instance.String())
		}
	}

	return output, nil
}

func formatHeading(text string) string {
	output := ""
	output += fmt.Sprintln()
	output += fmt.Sprintln(text)
	output += fmt.Sprintln("======================================")
	return output
}

func formatSubHeading(text string) string {
	output := ""
	output += fmt.Sprintln()
	output += fmt.Sprintln(text)
	output += fmt.Sprintln("--------------------------------------")
	return output
}
