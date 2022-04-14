package git_test

import (
	"errors"
	"testing"

	"github.com/gandarez/changelog-action/pkg/git"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClean(t *testing.T) {
	gc := git.NewGit("/path/to/repo")

	value, err := gc.Clean("'test'", nil)
	require.NoError(t, err)

	assert.Equal(t, "test", value)
}

func TestCleanErr(t *testing.T) {
	gc := git.NewGit("/path/to/repo")

	value, err := gc.Clean("'test'", errors.New("error"))
	require.Error(t, err)

	assert.Equal(t, "test", value)
	assert.EqualError(t, err, "error")
}

func TestLatestTag(t *testing.T) {
	gc := git.NewGit("/path/to/repo")
	gc.GitCmd = func(env map[string]string, args ...string) (string, error) {
		assert.Nil(t, env)
		assert.Equal(t, args, []string{"tag", "--points-at", "HEAD", "--sort", "-version:creatordate"})

		return "v0.2.17", nil
	}

	value := gc.LatestTagOrHash()

	assert.Equal(t, "v0.2.17", value)
}

func TestLatestTag_NoTagFound(t *testing.T) {
	var numCalls int

	gc := git.NewGit("/path/to/repo")
	gc.GitCmd = func(env map[string]string, args ...string) (string, error) {
		numCalls++

		assert.Nil(t, env)

		switch numCalls {
		case 1:
			assert.Equal(t, args, []string{"tag", "--points-at", "HEAD", "--sort", "-version:creatordate"})
		case 2:
			assert.Equal(t, args, []string{"describe", "--tags", "--abbrev=0"})
		case 3:
			assert.Equal(t, args, []string{"rev-parse", "HEAD"})
		}

		return "", nil
	}

	value := gc.LatestTagOrHash()

	assert.Empty(t, value)
}

func TestPreviousTag(t *testing.T) {
	gc := git.NewGit("/path/to/repo")
	gc.GitCmd = func(env map[string]string, args ...string) (string, error) {
		assert.Nil(t, env)
		assert.Equal(t, args, []string{"describe", "--tags", "--abbrev=0", "tags/v1.4.9^"})

		return "v1.4.8", nil
	}

	value, err := gc.PreviousTag("v1.4.9")
	require.NoError(t, err)

	assert.Equal(t, "v1.4.8", value)
}

func TestPreviousTagErr(t *testing.T) {
	var numCalls int

	gc := git.NewGit("/path/to/repo")
	gc.GitCmd = func(env map[string]string, args ...string) (string, error) {
		numCalls++

		assert.Nil(t, env)

		switch numCalls {
		case 1:
			assert.Equal(t, args, []string{"describe", "--tags", "--abbrev=0", "tags/v1.4.9^"})
		case 2:
			assert.Equal(t, args, []string{"rev-list", "--max-parents=0", "HEAD"})
		}

		return "", errors.New("error")
	}

	_, err := gc.PreviousTag("v1.4.9")

	assert.EqualError(t, err, "error")
}

func TestLog(t *testing.T) {
	gc := git.NewGit("/path/to/repo")
	gc.GitCmd = func(env map[string]string, args ...string) (string, error) {
		assert.Nil(t, env)
		assert.Equal(t, args, []string{"log", "--pretty=oneline", "--abbrev-commit", "--no-decorate", "--no-color", "tags/v1.2.3..tags/v1.3.0"})

		return "2b982db Add workflows\n5a359bb Fix logging", nil
	}

	value, err := gc.Log("tags/v1.2.3..tags/v1.3.0")
	require.NoError(t, err)

	assert.Equal(t, "2b982db Add workflows\n5a359bb Fix logging", value)
}

func TestLogErr(t *testing.T) {
	gc := git.NewGit("/path/to/repo")
	gc.GitCmd = func(env map[string]string, args ...string) (string, error) {
		assert.Nil(t, env)
		assert.Equal(t, args, []string{"log", "--pretty=oneline", "--abbrev-commit", "--no-decorate", "--no-color", "tags/v1.2.3..tags/v1.3.0"})

		return "", errors.New("error")
	}

	_, err := gc.Log("tags/v1.2.3..tags/v1.3.0")

	assert.EqualError(t, err, "error")
}
