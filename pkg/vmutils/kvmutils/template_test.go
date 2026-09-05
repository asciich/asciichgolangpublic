package kvmutils_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/vmutils/kvmutils"
)

const (
	testVmName     = "manualtest"
	testDiskPath   = "/storage/libvirt_disks/default_pool/manualtest"
	testMacAddress = "b8:27:eb:ac:27:74"
)

// newTestCreateOptions builds a KvmCreateVmOptions with valid values for the template rendering.
func newTestCreateOptions(t *testing.T) *kvmutils.KvmCreateVmOptions {
	t.Helper()

	createOptions := kvmutils.NewKvmCreateVmOptions()

	require.NoError(t, createOptions.SetVmName(testVmName))
	require.NoError(t, createOptions.SetDiskImageByPath(testDiskPath))
	require.NoError(t, createOptions.SetMacAddress(testMacAddress))

	return createOptions
}

// newTestCreateOptionsWithoutMac builds options WITHOUT a MAC address to prove it is optional.
func newTestCreateOptionsWithoutMac(t *testing.T) *kvmutils.KvmCreateVmOptions {
	t.Helper()

	createOptions := kvmutils.NewKvmCreateVmOptions()

	require.NoError(t, createOptions.SetVmName(testVmName))
	require.NoError(t, createOptions.SetDiskImageByPath(testDiskPath))
	// Intentionally NOT setting a MAC address.

	return createOptions
}

func TestLibvirtXmls_CreateXmlForVmAsString(t *testing.T) {
	t.Run("nil createOptions returns error", func(t *testing.T) {
		_, err := kvmutils.LibvirtXmls().CreateXmlForVmAsString(nil)
		require.Error(t, err)
	})

	t.Run("rendered xml contains the vm name", func(t *testing.T) {
		xmlString, err := kvmutils.LibvirtXmls().CreateXmlForVmAsString(newTestCreateOptions(t))
		require.NoError(t, err)
		require.Contains(t, xmlString, "<name>"+testVmName+"</name>")
	})

	t.Run("rendered xml contains the disk path", func(t *testing.T) {
		xmlString, err := kvmutils.LibvirtXmls().CreateXmlForVmAsString(newTestCreateOptions(t))
		require.NoError(t, err)
		require.Contains(t, xmlString, "source file='"+testDiskPath+"'")
	})

	t.Run("rendered xml contains the mac address", func(t *testing.T) {
		xmlString, err := kvmutils.LibvirtXmls().CreateXmlForVmAsString(newTestCreateOptions(t))
		require.NoError(t, err)
		require.Contains(t, xmlString, "mac address='"+testMacAddress+"'")
	})

	t.Run("rendered xml uses vnc graphics", func(t *testing.T) {
		xmlString, err := kvmutils.LibvirtXmls().CreateXmlForVmAsString(newTestCreateOptions(t))
		require.NoError(t, err)
		require.Contains(t, xmlString, "type='vnc'")
		require.NotContains(t, xmlString, "type='spice'")
	})

	t.Run("rendered xml has no unrendered template placeholders", func(t *testing.T) {
		xmlString, err := kvmutils.LibvirtXmls().CreateXmlForVmAsString(newTestCreateOptions(t))
		require.NoError(t, err)
		require.False(t, strings.Contains(xmlString, "{{"))
		require.False(t, strings.Contains(xmlString, "}}"))
	})

	t.Run("rendered xml is valid libvirt domain xml", func(t *testing.T) {
		xmlString, err := kvmutils.LibvirtXmls().CreateXmlForVmAsString(newTestCreateOptions(t))
		require.NoError(t, err)

		// GetMacAddressFromXmlString parses the XML via libvirt-go-xml, so a
		// successful parse proves the rendered document is well-formed.
		parsedMac, err := kvmutils.GetMacAddressFromXmlString(xmlString)
		require.NoError(t, err)
		require.EqualValues(t, testMacAddress, parsedMac)
	})
}

func TestLibvirtXmls_MacAddressIsOptional(t *testing.T) {
	t.Run("renders successfully without a mac address", func(t *testing.T) {
		xmlString, err := kvmutils.LibvirtXmls().CreateXmlForVmAsString(newTestCreateOptionsWithoutMac(t))
		require.NoError(t, err)
		require.NotEmpty(t, xmlString)
	})

	t.Run("rendered xml without mac contains no mac element", func(t *testing.T) {
		xmlString, err := kvmutils.LibvirtXmls().CreateXmlForVmAsString(newTestCreateOptionsWithoutMac(t))
		require.NoError(t, err)
		require.NotContains(t, xmlString, "<mac address=")
	})

	t.Run("rendered xml without mac still contains vm name and disk", func(t *testing.T) {
		xmlString, err := kvmutils.LibvirtXmls().CreateXmlForVmAsString(newTestCreateOptionsWithoutMac(t))
		require.NoError(t, err)
		require.Contains(t, xmlString, "<name>"+testVmName+"</name>")
		require.Contains(t, xmlString, "source file='"+testDiskPath+"'")
	})

	t.Run("rendered xml without mac still uses the network interface", func(t *testing.T) {
		xmlString, err := kvmutils.LibvirtXmls().CreateXmlForVmAsString(newTestCreateOptionsWithoutMac(t))
		require.NoError(t, err)
		require.Contains(t, xmlString, "<source network='default'/>")
		require.Contains(t, xmlString, "<model type='virtio'/>")
	})

	t.Run("rendered xml without mac has no unrendered placeholders", func(t *testing.T) {
		xmlString, err := kvmutils.LibvirtXmls().CreateXmlForVmAsString(newTestCreateOptionsWithoutMac(t))
		require.NoError(t, err)
		require.False(t, strings.Contains(xmlString, "{{"))
		require.False(t, strings.Contains(xmlString, "}}"))
	})

	t.Run("empty mac still renders (explicit empty string)", func(t *testing.T) {
		createOptions := kvmutils.NewKvmCreateVmOptions()
		require.NoError(t, createOptions.SetVmName(testVmName))
		require.NoError(t, createOptions.SetDiskImageByPath(testDiskPath))
		// MacAddress left as its zero value "".

		xmlString, err := kvmutils.LibvirtXmls().CreateXmlForVmAsString(createOptions)
		require.NoError(t, err)
		require.NotContains(t, xmlString, "<mac address=")
	})

	t.Run("invalid mac still returns error when provided", func(t *testing.T) {
		createOptions := kvmutils.NewKvmCreateVmOptions()
		require.NoError(t, createOptions.SetVmName(testVmName))
		require.NoError(t, createOptions.SetDiskImageByPath(testDiskPath))
		require.NoError(t, createOptions.SetMacAddress("not-a-valid-mac"))

		_, err := kvmutils.LibvirtXmls().CreateXmlForVmAsString(createOptions)
		require.Error(t, err)
	})

	t.Run("valid mac is still honored", func(t *testing.T) {
		xmlString, err := kvmutils.LibvirtXmls().CreateXmlForVmAsString(newTestCreateOptions(t))
		require.NoError(t, err)
		require.Contains(t, xmlString, "mac address='"+testMacAddress+"'")
	})
}

func TestLibvirtXmls_GetMacAddressFromXmlString(t *testing.T) {
	t.Run("empty string returns error", func(t *testing.T) {
		_, err := kvmutils.GetMacAddressFromXmlString("")
		require.Error(t, err)
	})

	t.Run("invalid xml returns error", func(t *testing.T) {
		_, err := kvmutils.GetMacAddressFromXmlString("this is not xml")
		require.Error(t, err)
	})

	t.Run("mac is extracted from rendered template", func(t *testing.T) {
		xmlString, err := kvmutils.LibvirtXmls().CreateXmlForVmAsString(newTestCreateOptions(t))
		require.NoError(t, err)

		parsedMac, err := kvmutils.GetMacAddressFromXmlString(xmlString)
		require.NoError(t, err)
		require.EqualValues(t, testMacAddress, parsedMac)
	})
}
