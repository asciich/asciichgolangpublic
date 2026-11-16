package containerimagehandler

import (
	"context"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// PushImage pushes the given image as 'nameAndTag' to a repository.
//
// The registry to push to can either be specified as part of 'nameAndTag'
// (e.g. 'registry.example.com/myimage:latest') or using 'options.Registry'.
// If both are set they must match, otherwise an error is returned.
//
// Additional fields in 'options' can be used to authenticate against the
// registry. If no credentials are set the default keychain is used.
//
// If 'nameAndTag' contains neither a tag nor a digest, the tag 'latest' is used.
//
// Returns the full image name and tag of the pushed image, including the
// registry if needed. The returned string can directly be used for e.g.
// 'docker run --rm -it <returned_string>'.
// In case of an error an empty string is returned.
func PushImage(ctx context.Context, image v1.Image, nameAndTag string, options *containeroptions.PushImageOptions) (string, error) {
	if image == nil {
		return "", tracederrors.TracedErrorNil("image")
	}

	if nameAndTag == "" {
		return "", tracederrors.TracedErrorEmptyString("nameAndTag")
	}

	if options == nil {
		options = &containeroptions.PushImageOptions{}
	}

	logging.LogInfoByCtxf(ctx, "Push container image as '%s' started.", nameAndTag)

	if !hasTagOrDigest(nameAndTag) {
		nameAndTag += ":latest"
		logging.LogInfoByCtxf(ctx, "Going to push as latest: '%s'", nameAndTag)
	}

	ret, err := mergeRegistryIntoNameAndTag(nameAndTag, options.Registry)
	if err != nil {
		return "", err
	}

	ref, err := name.ParseReference(ret)
	if err != nil {
		return "", tracederrors.TracedErrorf("Failed to parse reference '%s': %w", ret, err)
	}

	remoteOptions := []remote.Option{
		remoteAuthOption(options),
	}

	err = remote.Write(ref, image, remoteOptions...)
	if err != nil {
		return "", tracederrors.TracedErrorf("Failed to push container image as '%s': %w", ret, err)
	}

	logging.LogChangedByCtxf(ctx, "Pushed container image as '%s'.", ret)

	logging.LogInfoByCtxf(ctx, "Push container image as '%s' finished.", ret)

	return ret, nil
}

// remoteAuthOption builds the 'remote.Option' used to authenticate against
// the registry. If 'options.Username' is set, basic auth with the given
// credentials is used. Otherwise the default keychain (e.g. the local
// docker credential store) is used.
func remoteAuthOption(options *containeroptions.PushImageOptions) remote.Option {
	if options.Username != "" {
		return remote.WithAuth(&authn.Basic{
			Username: options.Username,
			Password: options.Password,
		})
	}

	return remote.WithAuthFromKeychain(authn.DefaultKeychain)
}

// hasTagOrDigest returns true if 'nameAndTag' contains a tag or a digest.
//
// Only the last path segment is checked for a tag, since a ':' in an
// earlier segment belongs to the registry port
// (e.g. 'localhost:5000/myimage' has no tag).
func hasTagOrDigest(nameAndTag string) bool {
	if strings.Contains(nameAndTag, "@") {
		return true
	}

	lastSegment := nameAndTag[strings.LastIndex(nameAndTag, "/")+1:]

	return strings.Contains(lastSegment, ":")
}

// hasExplicitRegistry returns true if 'nameAndTag' already specifies a
// registry (e.g. 'registry.example.com/myimage:latest' or
// 'localhost:5000/myimage:latest') instead of relying on the default
// registry (index.docker.io).
func hasExplicitRegistry(nameAndTag string) bool {
	if !strings.Contains(nameAndTag, "/") {
		return false
	}

	firstSegment := strings.SplitN(nameAndTag, "/", 2)[0]

	if strings.Contains(firstSegment, ".") {
		return true
	}

	if strings.Contains(firstSegment, ":") {
		return true
	}

	return firstSegment == "localhost"
}

// mergeRegistryIntoNameAndTag combines 'nameAndTag' and an optional
// 'registry' into a single reference string.
//
// If 'registry' is empty, 'nameAndTag' is returned unchanged.
// If 'nameAndTag' already contains an explicit registry and 'registry' is
// also set, both must match, otherwise an error is returned.
func mergeRegistryIntoNameAndTag(nameAndTag string, registry string) (string, error) {
	if registry == "" {
		return nameAndTag, nil
	}

	if hasExplicitRegistry(nameAndTag) {
		firstSegment := strings.SplitN(nameAndTag, "/", 2)[0]

		if firstSegment != registry {
			return "", tracederrors.TracedErrorf(
				"Conflicting registry specification: nameAndTag '%s' already specifies registry '%s' but options.Registry is set to '%s'.",
				nameAndTag, firstSegment, registry,
			)
		}

		return nameAndTag, nil
	}

	return registry + "/" + nameAndTag, nil
}
