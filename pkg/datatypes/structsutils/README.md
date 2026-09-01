# structsutils

Package `structsutils` provides utility functions for working with Go structs using reflection.

## Functions

- `IsStruct(objectToTest interface{}) bool` - Checks if the given object is a struct
- `IsPointerToStruct(objectToTest interface{}) bool` - Checks if the given object is a pointer to a struct
- `IsStructOrPointerToStruct(objectToTest interface{}) bool` - Checks if the given object is a struct or a pointer to a struct
- `GetFieldValueAsString(structToGetValueFrom any, fieldName string) (string, error)` - Gets the value of a struct field as a string
- `GetFieldValuesAsString(structToGetFieldsFrom any) ([]string, error)` - Gets all field values of a struct as strings
- `ListFieldNames(structToList any) ([]string, error)` - Lists all field names of a struct (sorted alphabetically)
- `HasField(structToTest any, fieldName string) (bool, error)` - Checks if a struct has a field with the given name

## Specifications

For specifications see [structsutils.spec.md](structsutils.spec.md)
