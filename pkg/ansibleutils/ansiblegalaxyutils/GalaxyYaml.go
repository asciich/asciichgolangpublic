package ansiblegalaxyutils

import (
	"context"

	"gitlab.asciich.ch/tools/asciichgolangpublic.git/files"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/contextutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/tracederrors"
	"gopkg.in/yaml.v3"
)

type GalaxyYaml struct {
	Namespace string   `yaml:"namespace"`
	Name      string   `yaml:"name"`
	Version   string   `yaml:"version"`
	Readme    string   `yaml:"readme"`
	Authors   []string `yaml:"authors"`
}

func WriteGalaxyYamlFromCreateCollectionOptions(ctx context.Context, outFile files.File, options *CreateCollectionFileStructureOptions) error {
	if outFile == nil {
		return tracederrors.TracedErrorNil("outFile")
	}

	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	namespace, err := options.GetNamespace()
	if err != nil {
		return err
	}

	name, err := options.GetName()
	if err != nil {
		return err
	}

	version, err := options.GetVersionAsString()
	if err != nil {
		return err
	}

	authors, err := options.GetAuthors()
	if err != nil {
		return err
	}

	return WriteGalaxyYaml(
		ctx,
		outFile,
		&GalaxyYaml{
			Namespace: namespace,
			Name:      name,
			Version:   version,
			Readme:    "README.md",
			Authors:   authors,
		},
	)
}

func WriteGalaxyYaml(ctx context.Context, outFile files.File, data *GalaxyYaml) error {
	if outFile == nil {
		return tracederrors.TracedErrorNil("outFile")
	}

	if data == nil {
		return tracederrors.TracedErrorNil("data")
	}

	toWrite, err := yaml.Marshal(data)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to marshal GalaxyYaml: %w", err)
	}

	return outFile.WriteBytes(toWrite, contextutils.GetVerboseFromContext(ctx))
}
