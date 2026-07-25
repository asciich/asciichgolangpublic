package openhandsutils

import (
	"context"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/versionutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

var (
	testOpenhandsOnce     sync.Once
	testOpenhandsInstance *Openhands
	testOpenhandsErr      error
)

func startTestOpenHands(ctx context.Context, t *testing.T) *Openhands {
	t.Helper()

	const port = 12800

	testOpenhandsOnce.Do(func() {
		_, testOpenhandsErr = StartAsDockerContainer(ctx, &StartContainerOptions{
			Port:                     port,
			ContainerName:            "openhands-testing",
			ReachableByOtherMachines: false,
			WorkspacePath:            ".",
		})
		if testOpenhandsErr != nil {
			return
		}

		testOpenhandsInstance, testOpenhandsErr = NewOpenHands("http://127.0.0.1:" + strconv.Itoa(port))
	})

	require.NoError(t, testOpenhandsErr)
	require.NotNil(t, testOpenhandsInstance)

	return testOpenhandsInstance
}

func TestNewOpenHands_ValidUrl(t *testing.T) {
	oh, err := NewOpenHands("http://localhost:8000")
	require.NoError(t, err)
	require.NotNil(t, oh)

	url, err := oh.GetUrl()
	assert.NoError(t, err)
	assert.Equal(t, "http://localhost:8000", url)
}

func TestNewOpenHands_InvalidUrl(t *testing.T) {
	oh, err := NewOpenHands("not-a-url")
	assert.Error(t, err)
	assert.Nil(t, oh)
}

func TestNewOpenHands_EmptyUrl(t *testing.T) {
	oh, err := NewOpenHands("")
	assert.Error(t, err)
	assert.Nil(t, oh)
}

func TestSetUrl_Valid(t *testing.T) {
	oh := &Openhands{}

	err := oh.SetUrl("http://localhost:8000")
	assert.NoError(t, err)
	assert.Equal(t, "http://localhost:8000", oh.Url)
}

func TestSetUrl_Invalid(t *testing.T) {
	oh := &Openhands{}

	err := oh.SetUrl("invalid")
	assert.Error(t, err)
	assert.Equal(t, "", oh.Url)
}

func TestSetUrl_Empty(t *testing.T) {
	oh := &Openhands{}

	err := oh.SetUrl("")
	assert.Error(t, err)
}

func TestGetUrl_Set(t *testing.T) {
	oh := &Openhands{Url: "http://localhost:8000"}

	url, err := oh.GetUrl()
	assert.NoError(t, err)
	assert.Equal(t, "http://localhost:8000", url)
}

func TestGetUrl_NotSet(t *testing.T) {
	oh := &Openhands{}

	url, err := oh.GetUrl()
	assert.Error(t, err)
	assert.Equal(t, "", url)
}

func TestSetUrl_OverwriteExisting(t *testing.T) {
	oh := &Openhands{Url: "http://localhost:8000"}

	err := oh.SetUrl("http://localhost:9000")
	assert.NoError(t, err)
	assert.Equal(t, "http://localhost:9000", oh.Url)
}

func TestNewOpenHands_HttpsUrl(t *testing.T) {
	oh, err := NewOpenHands("https://openhands.example.com")
	require.NoError(t, err)
	require.NotNil(t, oh)

	url, err := oh.GetUrl()
	assert.NoError(t, err)
	assert.Equal(t, "https://openhands.example.com", url)
}

func TestGetVersion(t *testing.T) {
	ctx := getCtx()
	openHands := startTestOpenHands(ctx, t)

	version, err := openHands.GetVersion(ctx)
	require.NoError(t, err)
	require.True(t, versionutils.IsSemanticVersionString(version))
}

func TestGetSessionApiKey(t *testing.T) {
	ctx := getCtx()
	openHands := startTestOpenHands(ctx, t)

	sessionApiKey, err := openHands.GetSessionApiKey(ctx)
	require.NoError(t, err)

	require.NotEmpty(t, sessionApiKey)
	require.Len(t, sessionApiKey, 64, "expected 64 char hex string")

	for _, c := range sessionApiKey {
		isHex := (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
		require.True(t, isHex, "expected hex character, got '%c'", c)
	}
}
