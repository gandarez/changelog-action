package main

import (
	"os"
	"strings"

	"github.com/gandarez/changelog-action/cmd/changelog"
	"github.com/gandarez/changelog-action/pkg/actions"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
)

func main() {
	log.SetHandler(cli.Default)

	result, err := changelog.Run()
	if err != nil {
		log.Errorf("failed to get changelog: %s\n", err)

		os.Exit(1)
	}

	outputFilepath := os.Getenv("GITHUB_OUTPUT")

	// Print changelog.
	log.Infof("CHANGELOG: %s", result)

	if err := actions.SetOutput(outputFilepath, "CHANGELOG", sanitize(result)); err != nil {
		log.Fatalf("%s\n", err)
	}
}

func sanitize(input string) string {
	input = strings.ReplaceAll(input, "%", "%25")
	input = strings.ReplaceAll(input, "\n", "%0A")
	input = strings.ReplaceAll(input, "\r", "%0D")

	return input
}
