package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aryann/difflib"
	"github.com/coreos/go-semver/semver"
	"github.com/fatih/color"
	"github.com/octoblu/docker-swarm-diff/helpers"
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
	expectationChan := helpers.GetExpectationsAsync()
	realityChan := helpers.GetRealityAsync()

	expectationResult := <-expectationChan
	fatalIfError("getExpectations", expectationResult.Err)
	realityResult := <-realityChan
	fatalIfError("getReality", realityResult.Err)

	formattedExpectation, err := helpers.Format(expectationResult.Servers)
	fatalIfError("format(expectation)", err)
	formattedReality, err := helpers.Format(realityResult.Servers)
	fatalIfError("format(reality)", err)

	diffs := difflib.Diff(
		strings.Split(formattedExpectation, "\n"),
		strings.Split(formattedReality, "\n"),
	)

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

func fatalIfError(msg string, err error) {
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
