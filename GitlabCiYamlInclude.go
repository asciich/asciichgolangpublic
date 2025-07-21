package asciichgolangpublic

import (
	"fmt"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type GitlabCiYamlInclude struct {
	Project string `yaml:"project"`
	File    string `yaml:"file"`
	Ref     string `yaml:"ref"`
}

func NewGitlabCiYamlInclude() (g *GitlabCiYamlInclude) {
	return new(GitlabCiYamlInclude)
}

// The UnmarshalYAML is used as custom unmarshal function to avoid issues with multinine file sections.
// E.g this works anyway:
//
// include:
//   - file: "c.yaml"
//
// But this is a problem without custom UnmarshalYaml function:
//
// include:
//   - file:
//   - "c.yaml"
func (g *GitlabCiYamlInclude) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	type SingleFileInclude struct {
		Project string `yaml:"project"`
		File    string `yaml:"file"`
		Ref     string `yaml:"ref"`
	}

	type MultiFileInclude struct {
		Project string   `yaml:"project"`
		File    []string `yaml:"file"`
		Ref     string   `yaml:"ref"`
	}

	singleFile := new(SingleFileInclude)
	multiFile := new(MultiFileInclude)

	err = unmarshal(singleFile)
	if err == nil {
		g.File = singleFile.File
		g.Project = singleFile.Project
		g.Ref = singleFile.Ref
	} else {
		err = unmarshal(multiFile)
		if err != nil {
			return tracederrors.TracedErrorf("Custom UnmarshalYAML for GitlabCiYamlInclude failed: %w", err)
		}

		if len(multiFile.File) <= 0 {
			return tracederrors.TracedErrorf("No files found")
		}

		if len(multiFile.File) != 1 {
			return tracederrors.TracedErrorf("Only one file supported at the moment but got '%v'", multiFile.File)
		}

		g.File = multiFile.File[0]
		g.Project = multiFile.Project
		g.Ref = multiFile.Ref
	}

	return nil
}

func (g *GitlabCiYamlInclude) EqualsIgnoreVersion(other *GitlabCiYamlInclude) (isEqual bool, err error) {
	if other == nil {
		return false, tracederrors.TracedError("other is nil")
	}

	thisProject, thisFile, err := g.GetProjectAndFile()
	if err != nil {
		return false, err
	}

	otherProject, otherFile, err := other.GetProjectAndFile()
	if err != nil {
		return false, err
	}

	if thisProject != otherProject {
		return false, nil
	}

	if thisFile != otherFile {
		return false, nil
	}

	return true, nil
}

func (g *GitlabCiYamlInclude) GetFile() (file string, err error) {
	if g.File == "" {
		return "", tracederrors.TracedErrorf("File not set")
	}

	return g.File, nil
}

func (g *GitlabCiYamlInclude) GetLoggableString() (loggableString string, err error) {
	project, file, err := g.GetProjectAndFile()
	if err != nil {
		return "", err
	}

	ref, err := g.GetRef()
	if err != nil {
		return "", err
	}

	loggableString = fmt.Sprintf(
		"gitlab ci yaml include: file '%s' of project '%s' in ref '%s'",
		file,
		project,
		ref,
	)

	return loggableString, nil
}

func (g *GitlabCiYamlInclude) GetProject() (project string, err error) {
	if g.Project == "" {
		return "", tracederrors.TracedErrorf("Project not set")
	}

	return g.Project, nil
}

func (g *GitlabCiYamlInclude) GetProjectAndFile() (project string, file string, err error) {
	project, err = g.GetProject()
	if err != nil {
		return "", "", err
	}

	file, err = g.GetFile()
	if err != nil {
		return "", "", err
	}

	return project, file, nil
}

func (g *GitlabCiYamlInclude) GetRef() (ref string, err error) {
	if g.Ref == "" {
		return "", tracederrors.TracedErrorf("Ref not set")
	}

	return g.Ref, nil
}

func (g *GitlabCiYamlInclude) IsEmpty() (isEmpty bool) {
	if g.Project != "" {
		return false
	}

	// TODO if g.File != "" {
	// TODO 	return false
	// TODO }

	if g.Ref != "" {
		return false
	}

	return true
}

func (g *GitlabCiYamlInclude) IsNonEmpty() (isNonEmpty bool) {
	return !g.IsEmpty()
}

func (g *GitlabCiYamlInclude) MustEqualsIgnoreVersion(other *GitlabCiYamlInclude) (isEqual bool) {
	isEqual, err := g.EqualsIgnoreVersion(other)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isEqual
}

func (g *GitlabCiYamlInclude) MustGetFile() (file string) {
	file, err := g.GetFile()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return file
}

func (g *GitlabCiYamlInclude) MustGetLoggableString() (loggableString string) {
	loggableString, err := g.GetLoggableString()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return loggableString
}

func (g *GitlabCiYamlInclude) MustGetProject() (project string) {
	project, err := g.GetProject()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return project
}

func (g *GitlabCiYamlInclude) MustGetProjectAndFile() (project string, file string) {
	project, file, err := g.GetProjectAndFile()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return project, file
}

func (g *GitlabCiYamlInclude) MustGetRef() (ref string) {
	ref, err := g.GetRef()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return ref
}

func (g *GitlabCiYamlInclude) MustSetFile(file string) {
	err := g.SetFile(file)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabCiYamlInclude) MustSetProject(project string) {
	err := g.SetProject(project)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabCiYamlInclude) MustSetRef(ref string) {
	err := g.SetRef(ref)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabCiYamlInclude) SetFile(file string) (err error) {
	if file == "" {
		return tracederrors.TracedErrorf("file is empty string")
	}

	g.File = file

	return nil
}

func (g *GitlabCiYamlInclude) SetProject(project string) (err error) {
	if project == "" {
		return tracederrors.TracedErrorf("project is empty string")
	}

	g.Project = project

	return nil
}

func (g *GitlabCiYamlInclude) SetRef(ref string) (err error) {
	if ref == "" {
		return tracederrors.TracedErrorf("ref is empty string")
	}

	g.Ref = ref

	return nil
}
