package changelog

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/apex/log"
	"github.com/gandarez/changelog-action/pkg/git"
)

type gitClient interface {
	IsRepo() bool
	MakeSafe() error
	LatestTagOrHash() string
	PreviousTag(tag string) (string, error)
	Log(refs ...string) (string, error)
}

func Run() (string, error) {
	params, err := LoadParams()
	if err != nil {
		return "", fmt.Errorf("failed to load parameters: %s", err)
	}

	if params.Debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("debug logs enabled\n")
	}

	log.Debug(params.String())

	git := git.NewGit(params.RepoDir)

	return Changelog(params, git)
}

func Changelog(params Params, gc gitClient) (string, error) {
	err := gc.MakeSafe()
	if err != nil {
		return "", fmt.Errorf("failed to make safe: %s", err)
	}

	if !gc.IsRepo() {
		return "", fmt.Errorf("current folder is not a git repository")
	}

	var tag = params.CurrentTag

	if tag == "" {
		tag = gc.LatestTagOrHash()
	}

	var refs = []string{fmt.Sprintf("%s..%s", params.PreviousTag, tag)}

	if params.PreviousTag == "" {
		previousTag, err := gc.PreviousTag(tag)
		if err != nil {
			return "", fmt.Errorf("failed to get previous tag: %s", err)
		}

		refs = []string{fmt.Sprintf("%s..%s", previousTag, tag)}
	}

	log, err := gc.Log(refs...)
	if err != nil {
		return "", fmt.Errorf("failed to get log: %s", err)
	}

	var entries = strings.Split(log, "\n")
	entries = entries[0 : len(entries)-1]

	entries, err = filterEntries(params.Exclude, entries)
	if err != nil {
		return "", err
	}

	changelogElements := []string{
		"## Changelog",
		strings.Join(entries, "\n"),
	}

	return strings.Join(changelogElements, "\n\n"), nil
}

func filterEntries(filters []string, entries []string) ([]string, error) {
	for _, filter := range filters {
		r, err := regexp.Compile(filter)
		if err != nil {
			return entries, err
		}

		entries = remove(r, entries)
	}

	return entries, nil
}

func remove(filter *regexp.Regexp, entries []string) []string {
	var result []string

	for _, entry := range entries {
		if !filter.MatchString(extractCommitInfo(entry)) {
			result = append(result, entry)
		}
	}

	return result
}

// extractCommitInfo removes first word which is the commit hash.
func extractCommitInfo(line string) string {
	return strings.Join(strings.Split(line, " ")[1:], " ")
}
