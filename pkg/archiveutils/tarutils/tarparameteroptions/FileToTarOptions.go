package tarparameteroptions

import "github.com/asciich/asciichgolangpublic/pkg/tracederrors"

type FileToTarOptions struct {
	// If set this file name is used in the archive instead:
	OverrideFileName string
}

func (f *FileToTarOptions) GetOverrideFileName() (string, error) {
	if f.OverrideFileName == "" {
		return "", tracederrors.TracedError("OverrideFileName not set")
	}

	return f.OverrideFileName, nil
}