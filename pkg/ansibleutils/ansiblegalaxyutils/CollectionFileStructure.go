package ansiblegalaxyutils

import (
	"context"

	"gitlab.asciich.ch/tools/asciichgolangpublic.git/files"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/contextutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/logging"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/tracederrors"
)

// Creates the files structure as documented on https://docs.ansible.com/ansible/latest/dev_guide/developing_collections_structure.html .
func CreateFileStructure(ctx context.Context, path string, options *CreateCollectionFileStructureOptions) error {
	if path == "" {
		return tracederrors.TracedErrorEmptyString("path")
	}

	dir, err := files.GetLocalDirectoryByPath(path)
	if err != nil {
		return err
	}

	return CreateFileStructureInDir(ctx, dir, options)
}

// Creates the files structure as documented on https://docs.ansible.com/ansible/latest/dev_guide/developing_collections_structure.html .
func CreateFileStructureInDir(ctx context.Context, dir files.Directory, options *CreateCollectionFileStructureOptions) error {
	if dir == nil {
		return tracederrors.TracedErrorNil("dir")
	}

	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	path, err := dir.GetPath()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Create ansible galaxy file structure in '%s' started.", path)

	_, err = dir.CreateSubDirectory("docx", contextutils.GetVerboseFromContext(ctx))
	if err != nil {
		return err
	}

	galaxyYamlFile, err := dir.CreateFileInDirectory(contextutils.GetVerboseFromContext(ctx), "galaxy.yml")
	if err != nil {
		return err
	}

	err = WriteGalaxyYamlFromCreateCollectionOptions(ctx, galaxyYamlFile, options)
	if err != nil {
		return err
	}

	metaDir, err := dir.CreateSubDirectory("meta", contextutils.GetVerboseFromContext(ctx))
	if err != nil {
		return err
	}

	_, err = metaDir.CreateFileInDirectory(contextutils.GetVerboseFromContext(ctx), "runtime.yml")
	if err != nil {
		return err
	}

	for _, d := range []string{"plugins", "roles", "playbooks", "tests"} {
		_, err := dir.CreateSubDirectory(d, contextutils.GetVerboseFromContext(ctx))
		if err != nil {
			return err
		}
	}

	readmeFile, err := dir.CreateFileInDirectory(contextutils.GetVerboseFromContext(ctx), "README.md")
	if err != nil {
		return err
	}

	err = WriteInitialCollectionReadmeFromCreateCollectionOptions(ctx, readmeFile, options)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Create ansible galaxy file structure in '%s' finished.", path)

	return nil
}
