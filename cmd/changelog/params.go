package changelog

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gandarez/changelog-action/pkg/actions"
)

type Params struct {
	CurrentTag  string
	PreviousTag string
	Exclude     []string
	RepoDir     string
	Debug       bool
}

func LoadParams() (Params, error) {
	var currentTag string

	if currentTagStr := actions.GetInput("current_tag"); currentTagStr != "" {
		currentTag = currentTagStr
	}

	var previousTag string

	if previousTagStr := actions.GetInput("previous_tag"); previousTagStr != "" {
		previousTag = previousTagStr
	}

	var exclude []string

	if excludeArr := actions.GetInput("exclude"); excludeArr != "" {
		exclude = strings.Split(excludeArr, "\n")
	}

	var repoDir = "."

	if repoDirStr := actions.GetInput("repo_dir"); repoDirStr != "" {
		repoDir = repoDirStr
	}

	var debug bool

	if debugStr := actions.GetInput("debug"); debugStr != "" {
		parsed, err := strconv.ParseBool(debugStr)
		if err != nil {
			return Params{}, fmt.Errorf("invalid debug argument: %s", debugStr)
		}

		debug = parsed
	}

	return Params{
		CurrentTag:  currentTag,
		PreviousTag: previousTag,
		Exclude:     exclude,
		RepoDir:     repoDir,
		Debug:       debug,
	}, nil
}

func (p Params) String() string {
	return fmt.Sprintf(
		"current tag: %q, previous tag: %q, exclude: %q, repo dir %q, debug: %t\n",
		p.CurrentTag,
		p.PreviousTag,
		strings.Join(p.Exclude, ","),
		p.RepoDir,
		p.Debug,
	)
}
