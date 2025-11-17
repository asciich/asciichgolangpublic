package atlassianconfluenceutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/wikiutils/atlassianconfluenceutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func TestGetPageIdFromUrl(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		pageId, err := atlassianconfluenceutils.GetPageIdFromUrl(getCtx(), "")
		require.Error(t, err)
		require.Empty(t, pageId)
	})

	t.Run("example url", func(t *testing.T) {
		pageId, err := atlassianconfluenceutils.GetPageIdFromUrl(getCtx(), "https://wiki.example.com/spaces/SPACE/pages/12345/Page+Title")
		require.NoError(t, err)
		require.EqualValues(t, "12345", pageId)
	})
}
