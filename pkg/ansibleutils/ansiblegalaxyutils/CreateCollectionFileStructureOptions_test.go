package ansiblegalaxyutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/ansibleutils/ansiblegalaxyutils"
)

func Test_GetVersionAsString(t *testing.T) {
	t.Run("Not set", func(t *testing.T) {
		options := &ansiblegalaxyutils.CreateCollectionFileStructureOptions{}
		versionString, err := options.GetVersionAsString()
		require.Error(t, err)
		require.Empty(t, versionString)
	})

	t.Run("empty string", func(t *testing.T) {
		options := &ansiblegalaxyutils.CreateCollectionFileStructureOptions{Version: ""}
		versionString, err := options.GetVersionAsString()
		require.Error(t, err)
		require.Empty(t, versionString)
	})

	t.Run("with leading v", func(t *testing.T) {
		options := &ansiblegalaxyutils.CreateCollectionFileStructureOptions{Version: "v0.1.2"}
		versionString, err := options.GetVersionAsString()
		require.NoError(t, err)
		require.EqualValues(t, "0.1.2", versionString) // Ansible galaxy expects without leading 'v'
	})

	t.Run("without leading v", func(t *testing.T) {
		options := &ansiblegalaxyutils.CreateCollectionFileStructureOptions{Version: "0.1.2"}
		versionString, err := options.GetVersionAsString()
		require.NoError(t, err)
		require.EqualValues(t, "0.1.2", versionString) // Ansible galaxy expects without leading 'v'
	})
}

func Test_GetName(t *testing.T) {
	t.Run("Not set", func(t *testing.T) {
		options := &ansiblegalaxyutils.CreateCollectionFileStructureOptions{}
		name, err := options.GetName()
		require.Error(t, err)
		require.Empty(t, name)
	})

	t.Run("Only one char", func(t *testing.T) {
		options := &ansiblegalaxyutils.CreateCollectionFileStructureOptions{Name: "a"}
		name, err := options.GetName()
		require.Error(t, err)
		require.Empty(t, name)
	})

	t.Run("Start with underscore", func(t *testing.T) {
		options := &ansiblegalaxyutils.CreateCollectionFileStructureOptions{Name: "_a"}
		name, err := options.GetName()
		require.Error(t, err)
		require.Empty(t, name)
	})

	t.Run("End with underscore", func(t *testing.T) {
		options := &ansiblegalaxyutils.CreateCollectionFileStructureOptions{Name: "a_"}
		name, err := options.GetName()
		require.Error(t, err)
		require.Empty(t, name)
	})

	t.Run("Include hyphen", func(t *testing.T) {
		options := &ansiblegalaxyutils.CreateCollectionFileStructureOptions{Name: "a-bc"}
		name, err := options.GetName()
		require.Error(t, err)
		require.Empty(t, name)
	})

	t.Run("Valid name", func(t *testing.T) {
		options := &ansiblegalaxyutils.CreateCollectionFileStructureOptions{Name: "abcdefg"}
		name, err := options.GetName()
		require.NoError(t, err)
		require.EqualValues(t, "abcdefg", name)
	})
}


func Test_GetNamespace(t *testing.T) {
	t.Run("Not set", func(t *testing.T) {
		options := &ansiblegalaxyutils.CreateCollectionFileStructureOptions{}
		namespace, err := options.GetNamespace()
		require.Error(t, err)
		require.Empty(t, namespace)
	})

	t.Run("Only one char", func(t *testing.T) {
		options := &ansiblegalaxyutils.CreateCollectionFileStructureOptions{Namespace: "a"}
		namespace, err := options.GetNamespace()
		require.Error(t, err)
		require.Empty(t, namespace)
	})

	t.Run("Start with underscore", func(t *testing.T) {
		options := &ansiblegalaxyutils.CreateCollectionFileStructureOptions{Namespace: "_a"}
		namespace, err := options.GetNamespace()
		require.Error(t, err)
		require.Empty(t, namespace)
	})

	t.Run("End with underscore", func(t *testing.T) {
		options := &ansiblegalaxyutils.CreateCollectionFileStructureOptions{Namespace: "a_"}
		namespace, err := options.GetNamespace()
		require.Error(t, err)
		require.Empty(t, namespace)
	})

	t.Run("Include hyphen", func(t *testing.T) {
		options := &ansiblegalaxyutils.CreateCollectionFileStructureOptions{Namespace: "a-bc"}
		namespace, err := options.GetNamespace()
		require.Error(t, err)
		require.Empty(t, namespace)
	})

	t.Run("Valid namespace", func(t *testing.T) {
		options := &ansiblegalaxyutils.CreateCollectionFileStructureOptions{Namespace: "abcdefg"}
		namespace, err := options.GetNamespace()
		require.NoError(t, err)
		require.EqualValues(t, "abcdefg", namespace)
	})
}

func Test_GetCreateFileOptions(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		collectionOptions := &ansiblegalaxyutils.CreateCollectionFileStructureOptions{}
		createOptions := collectionOptions.GetCreateFileOptions()
		require.False(t, createOptions.UseSudo)
	})
}