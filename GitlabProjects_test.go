package asciichgolangpublic

/* TODO enable again
import (
	"testing"


	"github.com/stretchr/testify/assert"
)


func TestGitlabProjectsProjectDoesNotExist(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"thisProjectDoesNotExist"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				gitlab := MustGetGitlabByFqdn("gitlab.asciich.ch")
				gitlab.MustUseUnauthenticatedClient(verbose)
				doesExist := gitlab.MustProjectByProjectPathExists("this/project_does_not_exist", verbose)
				assert.False(doesExist)
			},
		)
	}
}
*/