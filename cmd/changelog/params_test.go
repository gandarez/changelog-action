package changelog_test

import (
	"os"
	"testing"

	"github.com/gandarez/changelog-action/cmd/changelog"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadParams_CurrentTag(t *testing.T) {
	os.Setenv("INPUT_CURRENT_TAG", "v1.2.3")
	defer os.Unsetenv("INPUT_CURRENT_TAG")

	params, err := changelog.LoadParams()
	require.NoError(t, err)

	assert.Equal(t, "v1.2.3", params.CurrentTag)
}

func TestLoadParams_PreviousTag(t *testing.T) {
	os.Setenv("INPUT_PREVIOUS_TAG", "v0.2.3")
	defer os.Unsetenv("INPUT_PREVIOUS_TAG")

	params, err := changelog.LoadParams()
	require.NoError(t, err)

	assert.Equal(t, "v0.2.3", params.PreviousTag)
}

func TestLoadParams_Exclude(t *testing.T) {
	os.Setenv("INPUT_EXCLUDE", "^Merge .*\nFix .*")
	defer os.Unsetenv("INPUT_EXCLUDE")

	params, err := changelog.LoadParams()
	require.NoError(t, err)

	assert.Equal(t, []string{"^Merge .*", "Fix .*"}, params.Exclude)
}

func TestLoadParams_RepoDir(t *testing.T) {
	os.Setenv("INPUT_REPO_DIR", "/var/tmp/folder")
	defer os.Unsetenv("INPUT_REPO_DIR")

	params, err := changelog.LoadParams()
	require.NoError(t, err)

	assert.Equal(t, "/var/tmp/folder", params.RepoDir)
}

func TestLoadParams_Debug(t *testing.T) {
	os.Setenv("INPUT_DEBUG", "true")
	defer os.Unsetenv("INPUT_DEBUG")

	params, err := changelog.LoadParams()
	require.NoError(t, err)

	assert.True(t, params.Debug)
}

func TestLoadParams_DebugErr(t *testing.T) {
	os.Setenv("INPUT_DEBUG", "10")
	defer os.Unsetenv("INPUT_DEBUG")

	_, err := changelog.LoadParams()

	assert.Error(t, err)
}
