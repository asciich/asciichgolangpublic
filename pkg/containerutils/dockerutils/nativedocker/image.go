package nativedocker

import (
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type Image struct {
	name string
}

func NewImage() *Image {
	return new(Image)
}

func (i *Image) SetName(name string) error {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	i.name = name

	return nil
}

func (i *Image) GetName() (string, error) {
	if i.name == "" {
		return "", tracederrors.TracedError("name not set")
	}

	return i.name, nil
}

func (i *Image) Exists(ctx context.Context) (bool, error) {
	name, err := i.GetName()
	if err != nil {
		return false, err
	}

	return NewDocker().ImageExists(ctx, name)
}

func (i *Image) Remove(ctx context.Context, options *dockeroptions.RemoveOptions) error {
	name, err := i.GetName()
	if err != nil {
		return err
	}

	if options == nil {
		options = &dockeroptions.RemoveOptions{}
	}

	return NewDocker().RemoveImage(ctx, name, options)
}
