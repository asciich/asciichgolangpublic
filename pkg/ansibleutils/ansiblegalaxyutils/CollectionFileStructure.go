package ansiblegalaxyutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
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
func CreateFileStructureInDir(ctx context.Context, dir filesinterfaces.Directory, options *CreateCollectionFileStructureOptions) error {
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

	_, err = dir.CreateSubDirectory(ctx, "docx", &filesoptions.CreateOptions{})
	if err != nil {
		return err
	}

	createFilesOptions := options.GetCreateFileOptions()

	galaxyYamlFile, err := dir.CreateFileInDirectory(ctx, "galaxy.yml", createFilesOptions)
	if err != nil {
		return err
	}

	err = WriteGalaxyYamlFromCreateCollectionOptions(ctx, galaxyYamlFile, options)
	if err != nil {
		return err
	}

	metaDir, err := dir.CreateSubDirectory(ctx, "meta", createFilesOptions)
	if err != nil {
		return err
	}

	_, err = metaDir.CreateFileInDirectory(ctx, "runtime.yml", createFilesOptions)
	if err != nil {
		return err
	}

	for _, d := range []string{"plugins", "roles", "playbooks", "tests"} {
		_, err := dir.CreateSubDirectory(ctx, d, createFilesOptions)
		if err != nil {
			return err
		}
	}

	readmeFile, err := dir.CreateFileInDirectory(ctx, "README.md", createFilesOptions)
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
