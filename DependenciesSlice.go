package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/errors"
	"github.com/asciich/asciichgolangpublic/logging"
)

type DependenciesSliceService struct{}

func DependenciesSlice() (d *DependenciesSliceService) {
	return NewDependenciesSliceService()
}

func NewDependenciesSliceService() (d *DependenciesSliceService) {
	return new(DependenciesSliceService)
}

func (d *DependenciesSliceService) AddSourceFileForEveryEntry(dependencies []Dependency, sourceFile File) (err error) {
	if dependencies == nil {
		return errors.TracedErrorNil("dependencies")
	}

	if sourceFile == nil {
		return errors.TracedErrorNil("sourceFile")
	}

	for _, dependency := range dependencies {
		err = dependency.AddSourceFile(sourceFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *DependenciesSliceService) GetDependencyNames(dependencies []Dependency) (dependencyNames []string, err error) {
	for _, toAdd := range dependencies {
		nameToAdd, err := toAdd.GetName()
		if err != nil {
			return nil, err
		}

		dependencyNames = append(dependencyNames, nameToAdd)
	}

	return dependencyNames, nil
}

func (d *DependenciesSliceService) MustAddSourceFileForEveryEntry(dependencies []Dependency, sourceFile File) {
	err := d.AddSourceFileForEveryEntry(dependencies, sourceFile)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (d *DependenciesSliceService) MustGetDependencyNames(dependencies []Dependency) (dependencyNames []string) {
	dependencyNames, err := d.GetDependencyNames(dependencies)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return dependencyNames
}
