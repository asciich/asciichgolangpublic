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

func TestGetRunbookLinks_ReturnsSingleString(t *testing.T) {
	testCase := &testcase.TestCase{
		RunbookLinks: "https://example.com/runbook",
	}

	links, err := testCase.GetRunbookLinks()
	require.NoError(t, err)
	require.Len(t, links, 1)
	require.Equal(t, "https://example.com/runbook", links[0])
}

func TestGetRunbookLinks_ReturnsMultipleStrings(t *testing.T) {
	testCase := &testcase.TestCase{
		RunbookLinks: []any{"https://example.com/runbook1", "https://example.com/runbook2"},
	}

	links, err := testCase.GetRunbookLinks()
	require.NoError(t, err)
	require.Len(t, links, 2)
	require.Equal(t, "https://example.com/runbook1", links[0])
	require.Equal(t, "https://example.com/runbook2", links[1])
}

func TestGetRunbookLinks_ErrorWhenNil(t *testing.T) {
	testCase := &testcase.TestCase{
		RunbookLinks: nil,
	}

	links, err := testCase.GetRunbookLinks()
	require.Error(t, err)
	require.Nil(t, links)
}

func TestGetRunbookLinks_ErrorWhenEmptyString(t *testing.T) {
	testCase := &testcase.TestCase{
		RunbookLinks: "",
	}

	links, err := testCase.GetRunbookLinks()
	require.Error(t, err)
	require.Nil(t, links)
}

func TestGetRunbookLinks_ErrorWhenEmptySlice(t *testing.T) {
	testCase := &testcase.TestCase{
		RunbookLinks: []any{},
	}

	links, err := testCase.GetRunbookLinks()
	require.Error(t, err)
	require.Nil(t, links)
}

func TestGetHintsForInvestigation_ReturnsHints(t *testing.T) {
	testCase := &testcase.TestCase{
		HintsForInvestigation: "Check the logs for errors",
	}

	hints, err := testCase.GetHintsForInvestigation()
	require.NoError(t, err)
	require.Equal(t, "Check the logs for errors", hints)
}

func TestGetHintsForInvestigation_ErrorWhenEmpty(t *testing.T) {
	testCase := &testcase.TestCase{
		HintsForInvestigation: "",
	}

	hints, err := testCase.GetHintsForInvestigation()
	require.Error(t, err)
	require.Empty(t, hints)
}

func TestFormatFailedMessage_WithRunbookLinksAndHints(t *testing.T) {
	testCase := &testcase.TestCase{
		RunbookLinks:          []any{"https://example.com/runbook1", "https://example.com/runbook2"},
		HintsForInvestigation: "Check the logs",
	}

	message := testCase.FormatFailedMessage("Test failed")
	require.Contains(t, message, "Test failed")
	require.Contains(t, message, "Runbook links:")
	require.Contains(t, message, "https://example.com/runbook1")
	require.Contains(t, message, "https://example.com/runbook2")
	require.Contains(t, message, "Hints for investigation: Check the logs")
}

func TestFormatFailedMessage_WithoutRunbookLinksAndHints(t *testing.T) {
	testCase := &testcase.TestCase{
		RunbookLinks:          nil,
		HintsForInvestigation: "",
	}

	message := testCase.FormatFailedMessage("Test failed")
	require.Contains(t, message, "Test failed")
	require.Contains(t, message, "No runbook_links set.")
	require.Contains(t, message, "No hints_for_investigation set.")
}

func TestFormatFailedMessage_WithSingleRunbookLink(t *testing.T) {
	testCase := &testcase.TestCase{
		RunbookLinks:          "https://example.com/runbook",
		HintsForInvestigation: "Check the logs",
	}

	message := testCase.FormatFailedMessage("Test failed")
	require.Contains(t, message, "Test failed")
	require.Contains(t, message, "Runbook links:")
	require.Contains(t, message, "https://example.com/runbook")
	require.Contains(t, message, "Hints for investigation: Check the logs")
}
