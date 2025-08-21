package datatypes

import "github.com/asciich/asciichgolangpublic/pkg/datatypes/gettypename"

func GetTypeName(input interface{}) (typeName string, err error) {
	return gettypename.GetTypeName(input)
}
