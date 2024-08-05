package asciichgolangpublic

import (
	"testing"


	"github.com/stretchr/testify/assert"
)


func TestGitTagGetTagByHash(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				repo := TemporaryDirectories().MustCreateEmptyTemporaryGitRepository(
					&CreateRepositoryOptions{
						Verbose:                   verbose,
						BareRepository:            false,
						InitializeWithEmptyCommit: true,
						InitializeWithDefaultAuthor: true,
					})
				defer repo.Delete(verbose)

				currentCommitHash := repo.MustGetCurrentCommitHash()
				gitTag := repo.MustGetTagByHash(currentCommitHash)
				assert.EqualValues(
					currentCommitHash,
					gitTag.MustGetHash(),
				)
			},
		)
	}
}