package changelog

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilterEntries(t *testing.T) {
	filters := []string{
		"^Merge pull request .*",
		"Fix .*",
	}

	entries := []string{
		"2b982db Fix logging",
		"5a359bb Add git ignore",
		"55df180 Merge pull request #10 from author/bugfix/on_release",
	}

	filtered, err := filterEntries(filters, entries)
	require.NoError(t, err)

	assert.Equal(t, []string{
		"5a359bb Add git ignore",
	}, filtered)
}

func TestExtractCommitInfo(t *testing.T) {
	result := extractCommitInfo("55df180 Merge pull request #10 from author/bugfix/on_release")

	assert.Equal(t, "Merge pull request #10 from author/bugfix/on_release", result)
}
