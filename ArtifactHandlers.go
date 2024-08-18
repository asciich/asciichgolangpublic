package asciichgolangpublic


type ArtifactHandlersService struct{}

func ArtifactHandlers() (a *ArtifactHandlersService) {
	return NewArtifactHandlersService()
}

func NewArtifactHandlersService() (a *ArtifactHandlersService) {
	return new(ArtifactHandlersService)
}

func (a *ArtifactHandlersService) GetArtifactHandlerForArtifact(artifactHandlers []ArtifactHandler, artifactName string) (handler ArtifactHandler, err error) {
	if artifactName == "" {
		return nil, TracedErrorEmptyString("artifactName")
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

	return nil, TracedErrorf("No artifact handler for '%s' found", artifactName)
}

func (a *ArtifactHandlersService) MustGetArtifactHandlerForArtifact(artifactHandlers []ArtifactHandler, artifactName string) (handler ArtifactHandler) {
	handler, err := a.GetArtifactHandlerForArtifact(artifactHandlers, artifactName)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return handler
}
