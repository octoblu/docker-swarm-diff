package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aryann/difflib"
	"github.com/coreos/go-semver/semver"
	"github.com/fatih/color"
	"github.com/octoblu/docker-swarm-diff/reality"
	"github.com/octoblu/docker-swarm-diff/server"
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
	app.Description = fmt.Sprintf(
		"Compare what docker swarm thinks should be running against what is actually running\n   %v\n   %v",
		color.RedString("- docker swarm thinks this should be running on this machine, but it isn't"),
		color.YellowString("+ this is running on a machine, but docker doesn't know about it"),
	)
	app.Flags = []cli.Flag{}
	app.Run(os.Args)
}

func run(context *cli.Context) {
	expectation, err := swarm.GetServers()
	panicIfError("swarm.GetServers", err)
	reality, err := reality.GetServers()
	panicIfError("reality.GetServers", err)

	formattedExpectation := strings.Split(format(expectation), "\n")
	formattedReality := strings.Split(format(reality), "\n")

	diffs := difflib.Diff(formattedExpectation, formattedReality)

	exitCode := 0

	for _, diff := range diffs {
		if diff.Delta == difflib.LeftOnly {
			color.Red(diff.String())
			exitCode = 1
			continue
		}

		if diff.Delta == difflib.RightOnly {
			color.Yellow(diff.String())
			exitCode = 1
			continue
		}

		fmt.Println(diff.String())
	}

	os.Exit(exitCode)
}

func format(servers []server.Server) string {
	output := ""

	for _, serv := range servers {
		output += formatSubHeading(serv.String())

		instances, err := serv.ServiceInstances()
		panicIfError("serv.ServiceInstances", err)
		for _, instance := range instances {
			output += fmt.Sprintln(instance.String())
		}
	}

	return output
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

func panicIfError(msg string, err error) {
	if err == nil {
		return
	}

	log.Panicln(msg, err.Error())
}

func printFormatted(heading string, servers []server.Server) {
	fmt.Println(formatHeading(heading))
	fmt.Println(format(servers))
}

func version() string {
	version, err := semver.NewVersion(VERSION)
	if err != nil {
		errorMessage := fmt.Sprintf("Error with version number: %v", VERSION)
		log.Panicln(errorMessage, err.Error())
	}
	return version.String()
}
