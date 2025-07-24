package contextutils

import (
	"context"
	"errors"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type ChangeIndicator struct {
	// true indicates you received a change result.
	ChangedResult bool
}

var ErrorNoChangeIndicatorPresent = errors.New("no change indicator present in given contest")

// Get a child context containing an cache indicator which can indicate if a result bases on change data or not.
func WithChangeIndicator(ctx context.Context) context.Context {
	ret := ctx
	if ret == nil {
		ret = context.Background()
	}

	return context.WithValue(ret, ChangeIndicator{}, &ChangeIndicator{})
}

func IsChanged(ctx context.Context) bool {
	changed, err := IsChangedResult(ctx)
	if err != nil {
		return false
	}

	return changed
}

// true means the returned result bases on change data, false means no cache was involved.
// If ctx has no ChangeInidicator as Value it an error is returned.
func IsChangedResult(ctx context.Context) (bool, error) {
	if ctx == nil {
		return false, tracederrors.TracedErrorf("%w: ctx is nil", ErrorNoChangeIndicatorPresent)
	}

	value := ctx.Value(ChangeIndicator{})
	if value == nil {
		return false, tracederrors.TracedErrorf("%w: no ChangeIndicator present", ErrorNoChangeIndicatorPresent)
	}

	cacheIndicator, ok := value.(*ChangeIndicator)
	if !ok {
		return false, tracederrors.TracedErrorf("%w: unable to covert to ChangeIndicator", ErrorNoChangeIndicatorPresent)
	}

	return cacheIndicator.ChangedResult, nil
}

func IsChangeIndicatorPresent(ctx context.Context) bool {
	return ctx.Value(ChangeIndicator{}) != nil
}

// Use this to set the cache indication.
// Set the parameter 'change' to true to indiacte the given result bases on caching.
//
// If no cache indicator is present or ctx is nil the ctx stays untouched.
// This allows to use SetCahceIndicator wherever you want but it's only changing ctx if requested by a ctx with change indicator.
func SetChangeIndicator(ctx context.Context, change bool) {
	if ctx == nil {
		return
	}

	if !IsChangeIndicatorPresent(ctx) {
		return
	}

	value := ctx.Value(ChangeIndicator{})
	if value == nil {
		return
	}

	cacheIndicator, ok := value.(*ChangeIndicator)
	if !ok {
		return
	}

	cacheIndicator.ChangedResult = change
}
