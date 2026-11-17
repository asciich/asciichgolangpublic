package containerimagehandler_test

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/registry"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/random"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containerimagehandler"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containeroptions"
)

// newTestRegistry starts an in-memory container registry.
// The caller must call Close() on the returned server.
func newTestRegistry() *httptest.Server {
	return httptest.NewServer(
		registry.New(registry.Logger(log.New(io.Discard, "", 0))),
	)
}

// newTestRegistryWithBasicAuth starts an in-memory container registry which
// requires the given basic auth credentials for every request.
// The caller must call Close() on the returned server.
func newTestRegistryWithBasicAuth(username string, password string) *httptest.Server {
	reg := registry.New(registry.Logger(log.New(io.Discard, "", 0)))

	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, pass, ok := r.BasicAuth()
			if !ok || user != username || pass != password {
				w.Header().Set("WWW-Authenticate", `Basic realm="test"`)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			reg.ServeHTTP(w, r)
		}),
	)
}

// getRegistryHost returns the 'host:port' of the given test registry.
func getRegistryHost(server *httptest.Server) string {
	return strings.TrimPrefix(server.URL, "http://")
}

func newTestImage(t *testing.T) v1.Image {
	t.Helper()

	image, err := random.Image(1024, 1)
	require.NoError(t, err)

	return image
}

// requireImageInRegistry pulls 'reference' from the registry and checks that
// its digest equals the digest of 'expected'.
func requireImageInRegistry(t *testing.T, reference string, expected v1.Image) {
	t.Helper()

	ref, err := name.ParseReference(reference)
	require.NoError(t, err)

	pulled, err := remote.Image(ref)
	require.NoError(t, err)

	expectedDigest, err := expected.Digest()
	require.NoError(t, err)

	pulledDigest, err := pulled.Digest()
	require.NoError(t, err)

	require.Equal(t, expectedDigest, pulledDigest)
}

func Test_PushImage_InvalidInput(t *testing.T) {
	ctx := context.Background()

	t.Run("nil image", func(t *testing.T) {
		pushed, err := containerimagehandler.PushImage(ctx, nil, "registry.example.com/test/image:latest", &containeroptions.PushImageOptions{})
		require.Error(t, err)
		require.Empty(t, pushed)
	})

	t.Run("empty nameAndTag", func(t *testing.T) {
		pushed, err := containerimagehandler.PushImage(ctx, newTestImage(t), "", &containeroptions.PushImageOptions{})
		require.Error(t, err)
		require.Empty(t, pushed)
	})
}

func Test_PushImage_RegistryInNameAndTag(t *testing.T) {
	ctx := context.Background()

	server := newTestRegistry()
	defer server.Close()

	host := getRegistryHost(server)

	t.Run("with options", func(t *testing.T) {
		image := newTestImage(t)
		nameAndTag := host + "/test/with-options:v1"

		pushed, err := containerimagehandler.PushImage(ctx, image, nameAndTag, &containeroptions.PushImageOptions{})
		require.NoError(t, err)
		require.Equal(t, nameAndTag, pushed)

		requireImageInRegistry(t, pushed, image)
	})

	t.Run("with nil options", func(t *testing.T) {
		image := newTestImage(t)
		nameAndTag := host + "/test/with-nil-options:v1"

		pushed, err := containerimagehandler.PushImage(ctx, image, nameAndTag, nil)
		require.NoError(t, err)
		require.Equal(t, nameAndTag, pushed)

		requireImageInRegistry(t, pushed, image)
	})

	t.Run("without tag is pushed as latest", func(t *testing.T) {
		image := newTestImage(t)

		pushed, err := containerimagehandler.PushImage(ctx, image, host+"/test/without-tag", &containeroptions.PushImageOptions{})
		require.NoError(t, err)
		require.Equal(t, host+"/test/without-tag:latest", pushed)

		requireImageInRegistry(t, pushed, image)
	})

	t.Run("push is idempotent", func(t *testing.T) {
		image := newTestImage(t)
		nameAndTag := host + "/test/idempotent:v1"

		for range 2 {
			pushed, err := containerimagehandler.PushImage(ctx, image, nameAndTag, &containeroptions.PushImageOptions{})
			require.NoError(t, err)
			require.Equal(t, nameAndTag, pushed)
		}

		requireImageInRegistry(t, nameAndTag, image)
	})
}

func Test_PushImage_RegistryInOptions(t *testing.T) {
	ctx := context.Background()

	server := newTestRegistry()
	defer server.Close()

	host := getRegistryHost(server)

	t.Run("registry only in options", func(t *testing.T) {
		image := newTestImage(t)

		pushed, err := containerimagehandler.PushImage(ctx, image, "test/only-options:v1", &containeroptions.PushImageOptions{
			Registry: host,
		})
		require.NoError(t, err)
		require.Equal(t, host+"/test/only-options:v1", pushed, "returned name must include the registry")

		requireImageInRegistry(t, pushed, image)
	})

	t.Run("registry only in options without tag is pushed as latest", func(t *testing.T) {
		image := newTestImage(t)

		pushed, err := containerimagehandler.PushImage(ctx, image, "test/only-options-no-tag", &containeroptions.PushImageOptions{
			Registry: host,
		})
		require.NoError(t, err)
		require.Equal(t, host+"/test/only-options-no-tag:latest", pushed)

		requireImageInRegistry(t, pushed, image)
	})

	t.Run("registry in options and nameAndTag match", func(t *testing.T) {
		image := newTestImage(t)

		pushed, err := containerimagehandler.PushImage(ctx, image, host+"/test/matching:v1", &containeroptions.PushImageOptions{
			Registry: host,
		})
		require.NoError(t, err)
		require.Equal(t, host+"/test/matching:v1", pushed, "registry must not be added twice")

		requireImageInRegistry(t, pushed, image)
	})

	t.Run("registry in options and nameAndTag conflict", func(t *testing.T) {
		pushed, err := containerimagehandler.PushImage(ctx, newTestImage(t), host+"/test/conflict:v1", &containeroptions.PushImageOptions{
			Registry: "other-registry.example.com",
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "Conflicting registry specification")
		require.Empty(t, pushed)
	})
}

func Test_PushImage_Authentication(t *testing.T) {
	ctx := context.Background()

	const username = "testuser"
	const password = "testpassword"

	server := newTestRegistryWithBasicAuth(username, password)
	defer server.Close()

	host := getRegistryHost(server)

	t.Run("correct credentials", func(t *testing.T) {
		image := newTestImage(t)

		pushed, err := containerimagehandler.PushImage(ctx, image, host+"/test/auth:v1", &containeroptions.PushImageOptions{
			Username: username,
			Password: password,
		})
		require.NoError(t, err)
		require.Equal(t, host+"/test/auth:v1", pushed)
	})

	t.Run("correct credentials and registry in options", func(t *testing.T) {
		image := newTestImage(t)

		pushed, err := containerimagehandler.PushImage(ctx, image, "test/auth-options:v1", &containeroptions.PushImageOptions{
			Registry: host,
			Username: username,
			Password: password,
		})
		require.NoError(t, err)
		require.Equal(t, host+"/test/auth-options:v1", pushed)
	})

	t.Run("wrong password", func(t *testing.T) {
		pushed, err := containerimagehandler.PushImage(ctx, newTestImage(t), host+"/test/auth-wrong:v1", &containeroptions.PushImageOptions{
			Username: username,
			Password: "wrong-password",
		})
		require.Error(t, err)
		require.Empty(t, pushed)
	})
}
