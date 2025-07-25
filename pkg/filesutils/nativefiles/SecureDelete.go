package nativefiles

import (
	"context"

	"github.com/lu4p/shred"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// Overwrite the file before deleting it to prevent data restore.
// Like 'shred' on linux systems.
func SecureDelete(ctx context.Context, path string) error {
	if path == "" {
		return tracederrors.TracedErrorEmptyString(path)
	}

	if Exists(contextutils.WithSilent(ctx), path) {
		shredconf := shred.Conf{Times: 1, Zeros: true, Remove: true}
		err := shredconf.Path(path)
		if err != nil {
			return tracederrors.TracedErrorf("Secure delete of '%s' failed: %w", path, err)
		}
		logging.LogChangedByCtxf(ctx, "Securely deleted '%s'.", path)
	} else {
		logging.LogInfoByCtxf(ctx, "'%s' already absent. Skip secure delete.", path)
	}

	return nil
}
