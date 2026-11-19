package repoinfoutils_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/repositoryutils/repoinfoutils"
)

func ExampleGetGitlabPort() {
	ctx := contextutils.ContextSilent()

	// Create a temporary repository directory containing a '.repoinfo.yml'.
	repoDir := mustutils.Must(os.MkdirTemp("", "repoinfoutils"))
	defer os.RemoveAll(repoDir)

	content := `---
repository_info:
  gitlab_port: 2200
`
	mustutils.Must0(os.WriteFile(filepath.Join(repoDir, ".repoinfo.yml"), []byte(content), 0o644))

	// Returns an error if 'gitlab_port' is not set.
	port := mustutils.Must(repoinfoutils.GetGitlabPort(ctx, repoDir))

	fmt.Println(port)

	// Output:
	// 2200
}
