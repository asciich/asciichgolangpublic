package repoinfoutils_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/repositoryutils/repoinfoutils"
)

func ExampleGetGitlabPortOrDefaultIfUnset() {
	ctx := contextutils.ContextSilent()

	// Create a temporary repository directory containing a '.repoinfo.yml' without 'gitlab_port'.
	repoDir := mustutils.Must(os.MkdirTemp("", "repoinfoutils"))
	defer os.RemoveAll(repoDir)

	content := `---
repository_info:
  gitlab_instance: "gitlab.com"
`
	mustutils.Must0(os.WriteFile(filepath.Join(repoDir, ".repoinfo.yml"), []byte(content), 0o644))

	// Since 'gitlab_port' is not set the default port 2222 is returned.
	port := mustutils.Must(repoinfoutils.GetGitlabPortOrDefaultIfUnset(ctx, repoDir))

	fmt.Println(port)

	// Output:
	// 2222
}
