package gettypename_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/gettypename"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfileoo"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func TestGetTypeName_EmptyString(t *testing.T) {
	typeName, err := gettypename.GetTypeName("")
	require.NoError(t, err)
	require.EqualValues(t, "string", typeName)
}

func TestGetTypeName_EmptyStringPtr(t *testing.T) {
	input := ""
	typeName, err := gettypename.GetTypeName(&input)
	require.NoError(t, err)
	require.EqualValues(t, "*string", typeName)
}

func TestGetTypeName_HelloWorld(t *testing.T) {
	typeName, err := gettypename.GetTypeName("Hello world")
	require.NoError(t, err)
	require.EqualValues(t, "string", typeName)
}

func TestGetTypeName_HelloWorldPtr(t *testing.T) {
	input := "Hello world"
	typeName, err := gettypename.GetTypeName(&input)
	require.NoError(t, err)
	require.EqualValues(t, "*string", typeName)
}

func TestGetTypeName_Int(t *testing.T) {
	typeName, err := gettypename.GetTypeName(42)
	require.NoError(t, err)
	require.EqualValues(t, "int", typeName)
}

func TestGetTypeName_IntPtr(t *testing.T) {
	input := 42
	typeName, err := gettypename.GetTypeName(&input)
	require.NoError(t, err)
	require.EqualValues(t, "*int", typeName)
}

func TestGetTypeName_Int64(t *testing.T) {
	var input int64 = 123456789
	typeName, err := gettypename.GetTypeName(input)
	require.NoError(t, err)
	require.EqualValues(t, "int64", typeName)
}

func TestGetTypeName_Int64Ptr(t *testing.T) {
	var input int64 = 123456789
	typeName, err := gettypename.GetTypeName(&input)
	require.NoError(t, err)
	require.EqualValues(t, "*int64", typeName)
}

func TestGetTypeName_CommandExecutorooDirectory(t *testing.T) {
	exec := commandexecutorexecoo.Exec()
	input, err := commandexecutorfileoo.NewDirectory(exec, "/does_not_exist")
	require.NoError(t, err)

	typeName, err := gettypename.GetTypeName(input)
	require.NoError(t, err)
	require.EqualValues(t, "*commandexecutorfileoo.Directory", typeName)
}

func TestGetTypeName_FmtErrorf(t *testing.T) {
	input := fmt.Errorf("something went wrong")
	typeName, err := gettypename.GetTypeName(input)
	require.NoError(t, err)
	require.EqualValues(t, "error{message='something went wrong'}", typeName)
}

func TestGetTypeName_TracedError(t *testing.T) {
	input := tracederrors.TracedError("something went wrong")
	typeName, err := gettypename.GetTypeName(input)
	require.NoError(t, err)
	require.EqualValues(t, "TracedErrorType{message='something went wrong'}", typeName)
}
