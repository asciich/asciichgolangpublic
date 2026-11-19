package repoinfoutils_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/repositoryutils/repoinfoutils"
)

func ExampleGetGitlabPath() {
	ctx := contextutils.ContextSilent()

	// Create a temporary repository directory containing a '.repoinfo.yml'.
	repoDir := mustutils.Must(os.MkdirTemp("", "repoinfoutils"))
	defer os.RemoveAll(repoDir)

	// The '.git' suffix was added by mistake.
	content := `---
repository_info:
  gitlab_path: "tools/asciichgolangpublic.git"
`
	mustutils.Must0(os.WriteFile(filepath.Join(repoDir, ".repoinfo.yml"), []byte(content), 0o644))

	// The returned path never includes the '.git' suffix.
	gitlabPath := mustutils.Must(repoinfoutils.GetGitlabPath(ctx, repoDir))

	fmt.Println(gitlabPath)

	// Output:
	// tools/asciichgolangpublic
}
