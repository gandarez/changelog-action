package git

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/apex/log"
)

// Client is an struct to run git.
type Client struct {
	repoDir string
	GitCmd  func(env map[string]string, args ...string) (string, error)
}

// NewGit creates a new git instance.
func NewGit(repoDir string) *Client {
	return &Client{
		repoDir: repoDir,
		GitCmd:  gitCmdFn,
	}
}

// gitCmdFn runs a git command with the specified env vars and returns its output or errors.
func gitCmdFn(env map[string]string, args ...string) (string, error) {
	var extraArgs = []string{
		"-c", "log.showSignature=false",
	}
	args = append(extraArgs, args...)
	/* #nosec */
	var cmd = exec.Command("git", args...)

	if env != nil {
		cmd.Env = []string{}
		for k, v := range env {
			cmd.Env = append(cmd.Env, k+"="+v)
		}
	}

	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	log.WithField("args", args).Debug("running git")
	err := cmd.Run()

	log.WithField("stdout", stdout.String()).
		WithField("stderr", stderr.String()).
		Debug("git result")

	if err != nil {
		return "", errors.New(stderr.String())
	}

	return stdout.String(), nil
}

// Clean the output.
func (c *Client) Clean(output string, err error) (string, error) {
	output = strings.ReplaceAll(strings.Split(output, "\n")[0], "'", "")
	if err != nil {
		err = errors.New(strings.TrimSuffix(err.Error(), "\n"))
	}

	return output, err
}

// Run runs a git command and returns its output or errors.
func (c *Client) Run(args ...string) (string, error) {
	return c.GitCmd(nil, args...)
}

// MakeSafe adds safe.directory global config.
func (c *Client) MakeSafe() error {
	dir, err := filepath.Abs(c.repoDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for: %s", c.repoDir)
	}

	_, err = c.Run("config", "--global", "--add", "safe.directory", dir)
	if err != nil {
		return fmt.Errorf("failed to set safe current directory")
	}

	return nil
}

// IsRepo returns true if current folder is a git repository.
func (c *Client) IsRepo() bool {
	out, err := c.Run("rev-parse", "--is-inside-work-tree")
	return err == nil && strings.TrimSpace(out) == "true"
}

// LatestTag returns the latest tag or commit hash.
func (c *Client) LatestTagOrHash() string {
	for _, fn := range []func() (string, error){
		func() (string, error) {
			return c.Clean(c.Run("tag", "--points-at", "HEAD", "--sort", "-version:creatordate"))
		},
		func() (string, error) {
			return c.Clean(c.Run("describe", "--tags", "--abbrev=0"))
		},
		func() (string, error) {
			return c.Clean(c.Run("rev-parse", "HEAD"))
		},
	} {
		tag, _ := fn()
		if tag != "" {
			return tag
		}
	}

	return ""
}

// PreviousTag returns the previous tag from passed tag.
func (c *Client) PreviousTag(tag string) (result string, err error) {
	result, err = c.Clean(c.Run("describe", "--tags", "--abbrev=0", fmt.Sprintf("tags/%s^", tag)))
	if err != nil {
		result, err = c.Clean(c.Run("rev-list", "--max-parents=0", "HEAD"))
	}
	return
}

func (c *Client) Log(refs ...string) (string, error) {
	var args = []string{"log", "--pretty=oneline", "--abbrev-commit", "--no-decorate", "--no-color"}
	args = append(args, refs...)
	return c.Run(args...)
}
