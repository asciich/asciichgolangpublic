package containerimagehandler

import (
	"context"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func PushImage(ctx context.Context, image v1.Image, nameAndTag string) error {
	if image == nil {
		return tracederrors.TracedErrorNil("image")
	}

	if nameAndTag == "" {
		return tracederrors.TracedErrorEmptyString("nameAndTag")
	}

	logging.LogInfoByCtxf(ctx, "Push container image as '%s' started.", nameAndTag)

	if !strings.Contains(nameAndTag, ":") {
		nameAndTag += ":latest"
		logging.LogInfoByCtxf(ctx, "Going to push as latest: '%s'", nameAndTag)
	}

	ref, err := name.ParseReference(nameAndTag)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to parse reference '%s': %w", nameAndTag, err)
	}

	err = remote.Write(ref, image)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to push container image as '%s': %w", nameAndTag, err)
	}

	logging.LogChangedByCtxf(ctx, "Pushed container image as '%s'.", nameAndTag)

	logging.LogInfoByCtxf(ctx, "Push container image as '%s' finished.", nameAndTag)

	return nil
}
