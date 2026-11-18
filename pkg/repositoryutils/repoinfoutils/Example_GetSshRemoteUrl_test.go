package repoinfoutils_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/repositoryutils/repoinfoutils"
)

func ExampleGetSshRemoteUrl() {
	ctx := contextutils.ContextSilent()

	// Create a temporary repository directory containing a '.repoinfo.yml'.
	repoDir := mustutils.Must(os.MkdirTemp("", "repoinfoutils"))
	defer os.RemoveAll(repoDir)

	content := `---
repository_info:
  gitlab_instance: "gitlab.com"
  gitlab_path: "tools/asciichgolangpublic"
`
	mustutils.Must0(os.WriteFile(filepath.Join(repoDir, ".repoinfo.yml"), []byte(content), 0o644))

	// Builds 'ssh://git@<gitlab_instance>/<gitlab_path>.git'.
	url := mustutils.Must(repoinfoutils.GetSshRemoteUrl(ctx, repoDir))

	fmt.Println(url)

	// Output:
	// ssh://git@gitlab.com/tools/asciichgolangpublic.git
}
