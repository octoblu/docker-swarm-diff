package main

import (
	"fmt"
	"log"
	"os"

	"github.com/coreos/go-semver/semver"
	"github.com/octoblu/docker-swarm-diff/diff"
	"github.com/octoblu/docker-swarm-diff/reality"
	"github.com/octoblu/docker-swarm-diff/swarm"
	"github.com/urfave/cli"
	De "github.com/visionmedia/go-debug"
)

var debug = De.Debug("docker-swarm-diff:main")

func main() {
	app := cli.NewApp()
	app.Name = "docker-swarm-diff"
	app.Version = version()
	app.Action = run
	app.Flags = []cli.Flag{}
	app.Run(os.Args)
}

func run(context *cli.Context) {
	expectations, err := swarm.GetServers()
	panicIfError("swarm.GetServers", err)
	reality, err := reality.GetServers()
	panicIfError("reality.GetServers", err)

	diffs := diff.Differentiate(expectations, reality)
	if len(diffs) == 0 {
		os.Exit(0)
	}

	for _, diff := range diffs {
		fmt.Println(diff.String())
	}
	os.Exit(1)
}

func panicIfError(msg string, err error) {
	if err == nil {
		return
	}

	log.Panicln(msg, err.Error())
}

func version() string {
	version, err := semver.NewVersion(VERSION)
	if err != nil {
		errorMessage := fmt.Sprintf("Error with version number: %v", VERSION)
		log.Panicln(errorMessage, err.Error())
	}
	return version.String()
}
