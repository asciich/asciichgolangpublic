package pfsenseutils

import (
	"context"
	"io"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func (r *Router) GetRequest(ctx context.Context, url string) ([]byte, error) {
	if url == "" {
		return nil, tracederrors.TracedErrorEmptyString("url")
	}

	err := r.Login(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := r.client.Get(url)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to get pfSense router page '%s': %w", url, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to read response body of pfSense router page: %w", err)
	}

	return body, nil
}
