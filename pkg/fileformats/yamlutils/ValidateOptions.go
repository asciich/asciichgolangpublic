package yamlutils

type ValidateOptions struct {
	// Every JSON file by definition is a valid YAML.
	// If RefuesePureJson is set to true pure JSON files will not be accepted as valid YAML file.
	RefuesePureJson bool
}