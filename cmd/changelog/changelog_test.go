package changelog_test

import (
	"errors"
	"testing"

	"github.com/gandarez/changelog-action/cmd/changelog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChangelog(t *testing.T) {
	tests := map[string]struct {
		LatestTagOrHash string
		PreviousTag     string
		TagExists       bool
		Params          changelog.Params
		Expected        string
	}{
		"no previous tag": {
			PreviousTag: "e63c125b28842b17546cc92f635d7eccc8e909a7",
			Expected: "## Changelog\n\n" +
				"2b982db First commit",
		},
		"auto": {
			LatestTagOrHash: "v0.2.0",
			PreviousTag:     "v0.1.0",
			Expected: "## Changelog\n\n" +
				"2b982db First commit\n" +
				"5a359bb Second commit\n" +
				"1774db0 Merge pull request #1 from author/feature/feat-1",
		},
		"auto and latest tag is hash": {
			LatestTagOrHash: "e63c125b28842b17546cc92f635d7eccc8e909a7",
			PreviousTag:     "53db8447314a82e42e801568a085d424a739260a",
			Expected: "## Changelog\n\n" +
				"2b982db First commit\n" +
				"5a359bb Second commit\n" +
				"1774db0 Merge pull request #1 from author/feature/feat-1",
		},
		"current tag set": {
			LatestTagOrHash: "",
			PreviousTag:     "v0.2.0",
			TagExists:       true,
			Params: changelog.Params{
				CurrentTag: "v0.3.0",
			},
			Expected: "## Changelog\n\n" +
				"5a359bb Second commit\n" +
				"c57f56f Third commit",
		},
		"current tag and previous tag set": {
			LatestTagOrHash: "",
			PreviousTag:     "",
			TagExists:       true,
			Params: changelog.Params{
				CurrentTag:  "v0.3.0",
				PreviousTag: "v0.1.0",
			},
			Expected: "## Changelog\n\n" +
				"2b982db First commit\n" +
				"5a359bb Second commit\n" +
				"c57f56f Third commit",
		},
		"current tag set hash with previous tag set but does not exist": {
			LatestTagOrHash: "",
			PreviousTag:     "53db8447314a82e42e801568a085d424a739260a",
			TagExists:       false,
			Params: changelog.Params{
				CurrentTag:  "e63c125b28842b17546cc92f635d7eccc8e909a7",
				PreviousTag: "v0.0.0",
				Exclude: []string{
					"^Merge pull request .*",
				},
			},
			Expected: "## Changelog\n\n" +
				"2b982db First commit\n" +
				"5a359bb Second commit",
		},
		"auto and exclude": {
			LatestTagOrHash: "v0.2.0",
			PreviousTag:     "v0.1.0",
			Params: changelog.Params{
				Exclude: []string{
					"^Merge pull request .*",
				},
			},
			Expected: "## Changelog\n\n" +
				"2b982db First commit\n" +
				"5a359bb Second commit",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			gc := initGitClientMock(
				test.LatestTagOrHash,
				test.PreviousTag,
				test.TagExists,
			)

			result, err := changelog.Changelog(test.Params, gc)
			require.NoError(t, err)

			assert.Equal(t, test.Expected, result)
		})
	}
}

func TestChangelog_MakeSafeErr(t *testing.T) {
	gc := &gitClientMock{
		MakeSafeFn: func() error {
			return errors.New("error")
		},
	}

	_, err := changelog.Changelog(changelog.Params{}, gc)
	require.Error(t, err)

	assert.EqualError(t, err, "failed to make safe: error")
}

type gitClientMock struct {
	LatestTagOrHashFn        func() string
	LatestTagOrHashFnInvoked int
	IsRepoFn                 func() bool
	IsRepoFnInvoked          int
	MakeSafeFn               func() error
	MakeSafeFnInvoked        int
	PreviousTagFn            func(tag string) (string, error)
	PreviousTagFnInvoked     int
	TagExistsFn              func(tag string) bool
	TagExistsFnInvoked       int
	LogFn                    func(refs ...string) (string, error)
	LogFnInvoked             int
}

func initGitClientMock(latestTag, previousTag string, tagExists bool) *gitClientMock {
	return &gitClientMock{
		IsRepoFn: func() bool {
			return true
		},
		MakeSafeFn: func() error {
			return nil
		},
		LatestTagOrHashFn: func() string {
			return latestTag
		},
		PreviousTagFn: func(_ string) (string, error) {
			return previousTag, nil
		},
		TagExistsFn: func(_ string) bool {
			return tagExists
		},
		LogFn: func(refs ...string) (string, error) {
			switch refs[0] {
			case "e63c125b28842b17546cc92f635d7eccc8e909a7..":
				return "2b982db First commit\n", nil
			case "v0.1.0..v0.2.0":
				return "2b982db First commit\n" +
					"5a359bb Second commit\n" +
					"1774db0 Merge pull request #1 from author/feature/feat-1\n", nil
			case "53db8447314a82e42e801568a085d424a739260a..e63c125b28842b17546cc92f635d7eccc8e909a7":
				return "2b982db First commit\n" +
					"5a359bb Second commit\n" +
					"1774db0 Merge pull request #1 from author/feature/feat-1\n", nil
			case "v0.2.0..v0.3.0":
				return "5a359bb Second commit\n" +
					"c57f56f Third commit\n", nil
			case "v0.1.0..v0.3.0":
				return "2b982db First commit\n" +
					"5a359bb Second commit\n" +
					"c57f56f Third commit\n", nil
			default:
				return "", errors.New("no tag found")
			}
		},
	}
}

func (m *gitClientMock) IsRepo() bool {
	m.IsRepoFnInvoked++
	return m.IsRepoFn()
}

func (m *gitClientMock) MakeSafe() error {
	m.MakeSafeFnInvoked++
	return m.MakeSafeFn()
}

func (m *gitClientMock) LatestTagOrHash() string {
	m.LatestTagOrHashFnInvoked++
	return m.LatestTagOrHashFn()
}

func (m *gitClientMock) PreviousTag(tag string) (string, error) {
	m.PreviousTagFnInvoked++
	return m.PreviousTagFn(tag)
}

func (m *gitClientMock) TagExists(tag string) bool {
	m.TagExistsFnInvoked++
	return m.TagExistsFn(tag)
}

func (m *gitClientMock) Log(refs ...string) (string, error) {
	m.LogFnInvoked++
	return m.LogFn(refs...)
}
