package repoinfoutils_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/repositoryutils/repoinfoutils"
)

func ExampleGetGitlabInstance() {
	ctx := contextutils.ContextSilent()

	// Create a temporary repository directory containing a '.repoinfo.yml'.
	repoDir := mustutils.Must(os.MkdirTemp("", "repoinfoutils"))
	defer os.RemoveAll(repoDir)

	content := `---
repository_info:
  gitlab_instance: "gitlab.com"
`
	mustutils.Must0(os.WriteFile(filepath.Join(repoDir, ".repoinfo.yml"), []byte(content), 0o644))

	// The directory containing the '.repoinfo.yml' (or the path to the file itself) can be used.
	instance := mustutils.Must(repoinfoutils.GetGitlabInstance(ctx, repoDir))

	fmt.Println(instance)

	// Output:
	// gitlab.com
}
