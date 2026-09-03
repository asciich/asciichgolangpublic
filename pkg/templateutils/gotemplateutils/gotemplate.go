package gotemplateutils

import (
	"bytes"
	"context"
	"strings"
	"text/template"

	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func RenderTemplateFromFileAsString(ctx context.Context, inputFile filesinterfaces.File, variables map[string]interface{}) (rendered string, err error) {
	if inputFile == nil {
		return "", tracederrors.TracedError("inputFile is nil")
	}

	if variables == nil {
		return "", tracederrors.TracedError("variables is nil")
	}

	inputString, err := inputFile.ReadAsString(ctx)
	if err != nil {
		return "", err
	}

	rendered, err = RenderTemplateFromStringAsString(inputString, variables)
	if err != nil {
		return "", err
	}

	return rendered, nil
}

func RenderTemplateFromFilePathAsString(ctx context.Context, inputFilePath string, variables map[string]interface{}) (rendered string, err error) {
	if inputFilePath == "" {
		return "", tracederrors.TracedError("inputFilePath is empty string")
	}

	if variables == nil {
		return "", tracederrors.TracedError("variables is nil")
	}

	inputFile, err := files.GetLocalFileByPath(inputFilePath)
	if err != nil {
		return "", err
	}

	rendered, err = RenderTemplateFromFileAsString(ctx, inputFile, variables)
	if err != nil {
		return "", err
	}

	return rendered, nil
}

func RenderTemplateFromStringAsString(inputString string, variables map[string]interface{}) (rendered string, err error) {
	if inputString == "" {
		return "", tracederrors.TracedError("inputString is empty string")
	}

	if variables == nil {
		return "", tracederrors.TracedError("variables is nil")
	}

	funcMap := template.FuncMap{
		"ToUpper": strings.ToUpper,
		"ToLower": strings.ToLower,
	}

	templ, err := template.New("test").Option("missingkey=error").Funcs(funcMap).Parse(inputString)
	if err != nil {
		return "", tracederrors.TracedError(err.Error())
	}

	var renderWriter bytes.Buffer
	err = templ.Execute(&renderWriter, variables)
	if err != nil {
		return "", tracederrors.TracedError(err.Error())
	}

	rendered = renderWriter.String()
	return rendered, nil
}
