package mustutils

import "gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/logging"

// Used to apply the Must pattern to any function returning a value and an error:
//
//	Usage: x := Must(functionReturningOneValueAndAnError())
//	If there is no error the return value is returned and stored into 'x'.
//	If there is an error it will be logged using `logging.LogGoErrorFatal(err)` which aborts the execution.
func Must[T any](v T, err error) T {
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
	return v
}

// Same as `Must` but for functions only returning an error.
func Must0(err error) {
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
	return
}

// Same as `Must` but supporting 2 return values.
func Must2[T1 any, T2 any](v1 T1, v2 T2, err error) (T1, T2) {
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
	return v1, v2
}

// Same as `Must` but supporting 3 return values.
func Must3[T1 any, T2 any, T3 any](v1 T1, v2 T2, v3 T3, err error) (T1, T2, T3) {
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
	return v1, v2, v3
}

// Same as `Must` but supporting 4 return values.
func Must4[T1 any, T2 any, T3 any, T4 any](v1 T1, v2 T2, v3 T3, v4 T4, err error) (T1, T2, T3, T4) {
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
	return v1, v2, v3, v4
}
