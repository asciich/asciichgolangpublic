package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/pkg/dependencyutils/dependencyinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type DependenciesSliceService struct{}

func DependenciesSlice() (d *DependenciesSliceService) {
	return NewDependenciesSliceService()
}

func NewDependenciesSliceService() (d *DependenciesSliceService) {
	return new(DependenciesSliceService)
}

func (d *DependenciesSliceService) AddSourceFileForEveryEntry(dependencies []dependencyinterfaces.Dependency, sourceFile filesinterfaces.File) (err error) {
	if dependencies == nil {
		return tracederrors.TracedErrorNil("dependencies")
	}

	if sourceFile == nil {
		return tracederrors.TracedErrorNil("sourceFile")
	}

	for _, dependency := range dependencies {
		err = dependency.AddSourceFile(sourceFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *DependenciesSliceService) GetDependencyNames(dependencies []dependencyinterfaces.Dependency) (dependencyNames []string, err error) {
	for _, toAdd := range dependencies {
		nameToAdd, err := toAdd.GetName()
		if err != nil {
			return nil, err
		}

		dependencyNames = append(dependencyNames, nameToAdd)
	}

	return dependencyNames, nil
}
