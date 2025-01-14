package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/errors"
	"github.com/asciich/asciichgolangpublic/logging"
)

type ArtifactHandlersService struct{}

func ArtifactHandlers() (a *ArtifactHandlersService) {
	return NewArtifactHandlersService()
}

func NewArtifactHandlersService() (a *ArtifactHandlersService) {
	return new(ArtifactHandlersService)
}

func (a *ArtifactHandlersService) GetArtifactHandlerForArtifact(artifactHandlers []ArtifactHandler, artifactName string) (handler ArtifactHandler, err error) {
	if artifactName == "" {
		return nil, errors.TracedErrorEmptyString("artifactName")
	}

	for _, handler = range artifactHandlers {
		isHandling, err := handler.IsHandlingArtifactByName(artifactName)
		if err != nil {
			return nil, err
		}

		if isHandling {
			return handler, nil
		}
	}

	return nil, errors.TracedErrorf("No artifact handler for '%s' found", artifactName)
}

func (a *ArtifactHandlersService) MustGetArtifactHandlerForArtifact(artifactHandlers []ArtifactHandler, artifactName string) (handler ArtifactHandler) {
	handler, err := a.GetArtifactHandlerForArtifact(artifactHandlers, artifactName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return handler
}
