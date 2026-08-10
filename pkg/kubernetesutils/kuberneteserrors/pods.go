package kuberneteserrors

import "errors"

// Pod-related errors
var (
	ErrPodInFailedState     = errors.New("pod in failed state")
	ErrPodNotFound          = errors.New("pod not found")
	ErrPodDeleted           = errors.New("pod already deleted")
	ErrPodNotRunning        = errors.New("pod not in running state")
	ErrPodNotSucceeded      = errors.New("pod not in succeeded state")
	ErrPodTimeout           = errors.New("timeout waiting for pod")
	ErrContainerNotFound    = errors.New("container not found in pod")
	ErrUnableToRetrieveLogs = errors.New("unable to retrieve pod logs")
)

// IsPodInFailedState returns true if the error is ErrPodInFailedState
func IsPodInFailedState(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrPodInFailedState)
}

// IsPodNotFound returns true if the error is ErrPodNotFound
func IsPodNotFound(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrPodNotFound)
}

// IsPodDeleted returns true if the error is ErrPodDeleted
func IsPodDeleted(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrPodDeleted)
}

// IsPodNotRunning returns true if the error is ErrPodNotRunning
func IsPodNotRunning(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrPodNotRunning)
}

// IsPodNotSucceeded returns true if the error is ErrPodNotSucceeded
func IsPodNotSucceeded(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrPodNotSucceeded)
}

// IsPodTimeout returns true if the error is ErrPodTimeout
func IsPodTimeout(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrPodTimeout)
}

// IsContainerNotFound returns true if the error is ErrContainerNotFound
func IsContainerNotFound(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrContainerNotFound)
}

// IsUnableToRetrieveLogs returns true if the error is ErrUnableToRetrieveLogs
func IsUnableToRetrieveLogs(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrUnableToRetrieveLogs)
}
