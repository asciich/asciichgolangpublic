package fzfutils

import (
	"context"

	"github.com/koki-develop/go-fzf"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func RunFuzzySearch(ctx context.Context, inputData []byte, options *SearchOptions) ([]string, error) {
	if inputData == nil {
		return nil, tracederrors.TracedErrorNil("inputData")
	}

	if options == nil {
		options = &SearchOptions{}
	}

	items := stringsutils.SplitLines(string(inputData), true)
	items = slicesutils.RemoveEmptyStrings(items)

	if len(items) == 0 {
		return nil, tracederrors.TracedError("no input provided")
	}

	opts := []fzf.Option{
		fzf.WithPrompt("> "),
		fzf.WithInputPlaceholder("Search..."),
	}
	if options.Multi {
		opts = append(opts, fzf.WithNoLimit(true)) // unlimited multi-select
	}

	f, err := fzf.New(opts...)
	if err != nil {
		return nil, tracederrors.TracedErrorf("failed to initialize fzf: %w", err)
	}

	idxs, err := f.Find(items, func(i int) string {
		return items[i]
	})
	if err != nil {
		return nil, err
	}

	result := []string{}
	for _, idx := range idxs {
		result = append(result, items[idx])
	}

	return result, nil
}
