package testcase_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testcase"
)

func TestGetName_ReturnsName(t *testing.T) {
	tc := &testcase.TestCase{Name: "my-test"}
	name, err := tc.GetName()
	require.NoError(t, err)
	assert.Equal(t, "my-test", name)
}

func TestGetName_ErrorWhenEmpty(t *testing.T) {
	tc := &testcase.TestCase{}
	_, err := tc.GetName()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name not set")
}

func TestGetTestType_ReturnsTestType(t *testing.T) {
	tc := &testcase.TestCase{TestType: "integration"}
	testType, err := tc.GetTestType()
	require.NoError(t, err)
	assert.Equal(t, "integration", testType)
}

func TestGetTestType_ErrorWhenEmpty(t *testing.T) {
	tc := &testcase.TestCase{}
	_, err := tc.GetTestType()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "test type not set")
}

func TestGetCommand_ReturnsCommand(t *testing.T) {
	tc := &testcase.TestCase{Command: "echo hello"}
	cmd, err := tc.GetCommand()
	require.NoError(t, err)
	assert.Equal(t, "echo hello", cmd)
}

func TestGetCommand_ErrorWhenEmpty(t *testing.T) {
	tc := &testcase.TestCase{}
	_, err := tc.GetCommand()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "command not set")
}

func TestGetHost_ReturnsHost(t *testing.T) {
	tc := &testcase.TestCase{Host: "localhost"}
	host, err := tc.GetHost()
	require.NoError(t, err)
	assert.Equal(t, "localhost", host)
}

func TestGetHost_ReturnsHostWithDomain(t *testing.T) {
	tc := &testcase.TestCase{Host: "example.com"}
	host, err := tc.GetHost()
	require.NoError(t, err)
	assert.Equal(t, "example.com", host)
}

func TestGetHost_ReturnsHostWithIP(t *testing.T) {
	tc := &testcase.TestCase{Host: "192.168.1.1"}
	host, err := tc.GetHost()
	require.NoError(t, err)
	assert.Equal(t, "192.168.1.1", host)
}

func TestGetHost_ErrorWhenEmpty(t *testing.T) {
	tc := &testcase.TestCase{}
	_, err := tc.GetHost()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "host not set")
}

func TestGetPort_ReturnsPort(t *testing.T) {
	tc := &testcase.TestCase{Port: "8080"}
	port, err := tc.GetPort()
	require.NoError(t, err)
	assert.Equal(t, 8080, port)
}

func TestGetPort_ErrorWhenEmpty(t *testing.T) {
	tc := &testcase.TestCase{}
	_, err := tc.GetPort()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "port not set")
}

func TestGetPort_ErrorWhenNotANumber(t *testing.T) {
	tc := &testcase.TestCase{Port: "abc"}
	_, err := tc.GetPort()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Failed to convert")
}

func TestGetPort_ErrorWhenFloatingPoint(t *testing.T) {
	tc := &testcase.TestCase{Port: "80.5"}
	_, err := tc.GetPort()
	assert.Error(t, err)
}

func TestGetNamespace_ReturnsNamespace(t *testing.T) {
	tc := &testcase.TestCase{Namespace: "default"}
	ns, err := tc.GetNamespace()
	require.NoError(t, err)
	assert.Equal(t, "default", ns)
}

func TestGetNamespace_ErrorWhenEmpty(t *testing.T) {
	tc := &testcase.TestCase{}
	_, err := tc.GetNamespace()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "namespace not set")
}

func TestGetCluster_ReturnsCluster(t *testing.T) {
	tc := &testcase.TestCase{Cluster: "prod-cluster"}
	cluster, err := tc.GetCluster()
	require.NoError(t, err)
	assert.Equal(t, "prod-cluster", cluster)
}

func TestGetCluster_ErrorWhenEmpty(t *testing.T) {
	tc := &testcase.TestCase{}
	_, err := tc.GetCluster()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cluster not set")
}

func TestGetResourceName_ReturnsResourceName(t *testing.T) {
	tc := &testcase.TestCase{ResourceName: "my-deployment"}
	rn, err := tc.GetResourceName()
	require.NoError(t, err)
	assert.Equal(t, "my-deployment", rn)
}

func TestGetResourceName_ErrorWhenEmpty(t *testing.T) {
	tc := &testcase.TestCase{}
	_, err := tc.GetResourceName()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "resource_name not set")
}

func TestSetData_SetsData(t *testing.T) {
	tc := &testcase.TestCase{}
	data := map[string]string{"key": "value"}
	err := tc.SetData(data)
	require.NoError(t, err)
}

func TestSetData_ErrorWhenNil(t *testing.T) {
	tc := &testcase.TestCase{}
	err := tc.SetData(nil)
	assert.Error(t, err)
}
