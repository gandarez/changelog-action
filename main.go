package main

import (
	"fmt"
	"os"

	"github.com/gandarez/changelog-action/cmd/changelog"

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

	// Print changelog.
	log.Debugf("CHANGELOG: %s", result)
	fmt.Printf("::set-output name=CHANGELOG::%s\n", result)

	os.Exit(0)
}
