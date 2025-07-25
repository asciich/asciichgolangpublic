package ansiblegalaxyutils

import (
	"context"
	"fmt"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/documentutils/markdowndocument"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func WriteInitialCollectionReadme(ctx context.Context, outFile filesinterfaces.File, name string, namespace string) error {
	if outFile == nil {
		return tracederrors.TracedErrorNil("outFile")
	}

	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	if name == "" {
		return tracederrors.TracedErrorEmptyString("namespace")
	}

	doc := markdowndocument.NewMarkDownDocument()

	err := doc.AddTitleByString(fmt.Sprintf("Ansible collection %s.%s", name, namespace))
	if err != nil {
		return err
	}

	content, err := doc.RenderAsString()
	if err != nil {
		return err
	}

	return outFile.WriteString(content, contextutils.GetVerboseFromContext(ctx))
}

func WriteInitialCollectionReadmeFromCreateCollectionOptions(ctx context.Context, outFile filesinterfaces.File, options *CreateCollectionFileStructureOptions) error {
	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	name, err := options.GetName()
	if err != nil {
		return err
	}

	namespace, err := options.GetNamespace()
	if err != nil {
		return err
	}

	return WriteInitialCollectionReadme(ctx, outFile, name, namespace)
}
