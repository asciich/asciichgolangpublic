package httpgeneric

import "context"

type progressEveryNBytes struct{}

func WithDownloadProgressEveryNMBytes(ctx context.Context, nMBytes int) context.Context {
	return WithDownloadProgressEveryNkBytes(ctx, 1024*nMBytes)
}

func WithDownloadProgressEveryNkBytes(ctx context.Context, nkBytes int) context.Context {
	return WithDownloadProgressEveryNBytes(ctx, 1024*nkBytes)
}

func WithDownloadProgressEveryNBytes(ctx context.Context, nBytes int) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if nBytes < 0 {
		nBytes = 0
	}

	return context.WithValue(ctx, progressEveryNBytes{}, nBytes)
}

func GetProgressEveryNBytes(ctx context.Context) int {
	if ctx == nil {
		return 0
	}

	val := ctx.Value(progressEveryNBytes{})
	if val == nil {
		return 0
	}

	ret, ok := val.(int)
	if !ok {
		return 0
	}

	return ret
}
