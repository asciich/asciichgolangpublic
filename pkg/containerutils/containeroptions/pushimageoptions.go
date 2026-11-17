package containeroptions

// PushImageOptions specifies how an image is pushed to a registry by
// containerimagehandler.PushImage.
type PushImageOptions struct {
	// Optional registry to push to.
	// Can be left empty if the registry is already part of the
	// 'nameAndTag' passed to containerimagehandler.PushImage.
	Registry string

	// Optional credentials to authenticate against the registry.
	// If Username is empty the default keychain (e.g. docker
	// credential store, '~/.docker/config.json' or environment based
	// credential helpers) is used instead.
	Username string
	Password string
}
