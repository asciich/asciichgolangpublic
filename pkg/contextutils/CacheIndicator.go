package contextutils

import (
	"context"
	"errors"

	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type CacheIndicator struct {
	// true indicates you received a cached result.
	CachedResult bool
}

var ErrorNoCacheIndicatorPresent = errors.New("no cache indicator present in given contest")

// Get a child context containing an cache indicator which can indicate if a result bases on cached data or not.
func WithCacheIndicator(ctx context.Context) context.Context {
	ret := ctx
	if ret == nil {
		ret = context.Background()
	}

	return context.WithValue(ret, CacheIndicator{}, &CacheIndicator{})
}

// true means the returned result bases on cached data, false means no cache was involved.
// If ctx has no CacheInidicator as Value it an error is returned.
func IsCachedResult(ctx context.Context) (bool, error) {
	if ctx == nil {
		return false, tracederrors.TracedErrorf("%w: ctx is nil", ErrorNoCacheIndicatorPresent)
	}

	value := ctx.Value(CacheIndicator{})
	if value == nil {
		return false, tracederrors.TracedErrorf("%w: no CacheIndicator present", ErrorNoCacheIndicatorPresent)
	}

	cacheIndicator, ok := value.(*CacheIndicator)
	if !ok {
		return false, tracederrors.TracedErrorf("%w: unable to covert to CacheIndicator", ErrorNoCacheIndicatorPresent)
	}

	return cacheIndicator.CachedResult, nil
}

func IsCacheIndicatorPresent(ctx context.Context) bool {
	return ctx.Value(CacheIndicator{}) != nil
}

// Use this to set the cache indication.
// Set the parameter 'cached' to true to indiacte the given result bases on caching.
//
// If no cache indicator is present or ctx is nil the ctx stays untouched.
// This allows to use SetCahceIndicator wherever you want but it's only changing ctx if requested by a ctx with change indicator.
func SetCacheIndicator(ctx context.Context, cached bool) {
	if ctx == nil {
		return
	}

	if !IsCacheIndicatorPresent(ctx) {
		return
	}

	value := ctx.Value(CacheIndicator{})
	if value == nil {
		return
	}

	cacheIndicator, ok := value.(*CacheIndicator)
	if !ok {
		return
	}

	cacheIndicator.CachedResult = cached
}
