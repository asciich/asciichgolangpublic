package dockergeneric_test

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockergeneric"
)

// buildTar creates a plain tar stream with the given regular files.
func buildTar(t *testing.T, files map[string]string) []byte {
	t.Helper()

	var ret bytes.Buffer
	tw := tar.NewWriter(&ret)
	for name, content := range files {
		require.NoError(t, tw.WriteHeader(&tar.Header{
			Typeflag: tar.TypeReg,
			Name:     "./" + name,
			Mode:     0o644,
			Size:     int64(len(content)),
		}))
		_, err := tw.Write([]byte(content))
		require.NoError(t, err)
	}
	require.NoError(t, tw.Close())

	return ret.Bytes()
}

// readTarEntries returns the content of all regular files and the names of all directories.
// All names are normalized.
func readTarEntries(t *testing.T, data []byte) (map[string]string, []string) {
	t.Helper()

	files := map[string]string{}
	dirs := []string{}

	tr := tar.NewReader(bytes.NewReader(data))
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			return files, dirs
		}
		require.NoError(t, err)

		name := dockergeneric.NormalizeEntryName(hdr.Name)
		switch hdr.Typeflag {
		case tar.TypeReg:
			content, err := io.ReadAll(tr)
			require.NoError(t, err)
			files[name] = string(content)
		case tar.TypeDir:
			dirs = append(dirs, name)
		}
	}
}

// writeMultiArchive creates a multi volume archive from the given volume contents.
func writeMultiArchive(t *testing.T, volumeNames []string, volumes map[string]map[string]string) []byte {
	t.Helper()

	var ret bytes.Buffer
	err := dockergeneric.WriteMultiVolumeArchive(
		context.Background(),
		&ret,
		volumeNames,
		func(volumeName string) (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(buildTar(t, volumes[volumeName]))), nil
		},
	)
	require.NoError(t, err)

	return ret.Bytes()
}

func Test_CheckVolumeName(t *testing.T) {
	for _, name := range []string{"a", "A1", "my-volume_1.x"} {
		t.Run("valid name '"+name+"'", func(t *testing.T) {
			require.NoError(t, dockergeneric.CheckVolumeName(name))
		})
	}

	for _, name := range []string{"", "../etc", "a/b", "-a", ".a", "a b", "a;rm", "a\nb"} {
		t.Run("invalid name '"+name+"'", func(t *testing.T) {
			require.Error(t, dockergeneric.CheckVolumeName(name))
		})
	}
}

func Test_NormalizeEntryName(t *testing.T) {
	testCases := map[string]string{
		"":          "",
		".":         "",
		"./":        "",
		"/":         "",
		"a":         "a",
		"./a":       "a",
		"/a":        "a",
		"./a/b/":    "a/b",
		"a/../../b": "b",
	}

	for input, expected := range testCases {
		t.Run("'"+input+"'", func(t *testing.T) {
			require.Equal(t, expected, dockergeneric.NormalizeEntryName(input))
		})
	}
}

func Test_GetVolumeDirectory(t *testing.T) {
	t.Run("directory is below volumes", func(t *testing.T) {
		require.Equal(t, "volumes/my-volume", dockergeneric.GetVolumeDirectory("my-volume"))
	})
}

func Test_Manifest(t *testing.T) {
	t.Run("marshal and parse round trip", func(t *testing.T) {
		content, err := dockergeneric.MarshalManifest([]string{"a", "b"})
		require.NoError(t, err)

		manifest, err := dockergeneric.ParseManifest(content)
		require.NoError(t, err)
		require.Equal(t, dockergeneric.ManifestVersion, manifest.Version)
		require.Equal(t, []string{"a", "b"}, manifest.Volumes)
	})

	t.Run("marshal uses the documented format", func(t *testing.T) {
		content, err := dockergeneric.MarshalManifest([]string{"a", "b"})
		require.NoError(t, err)
		require.JSONEq(t, `{"version":1,"volumes":["a","b"]}`, string(content))
	})

	t.Run("invalid json is rejected", func(t *testing.T) {
		_, err := dockergeneric.ParseManifest([]byte("not json"))
		require.Error(t, err)
	})

	t.Run("unsupported version is rejected", func(t *testing.T) {
		_, err := dockergeneric.ParseManifest([]byte(`{"version":99,"volumes":["a"]}`))
		require.Error(t, err)
	})

	t.Run("manifest without volumes is rejected", func(t *testing.T) {
		_, err := dockergeneric.ParseManifest([]byte(`{"version":1,"volumes":[]}`))
		require.Error(t, err)
	})

	t.Run("manifest with invalid volume name is rejected", func(t *testing.T) {
		_, err := dockergeneric.ParseManifest([]byte(`{"version":1,"volumes":["../etc"]}`))
		require.Error(t, err)
	})
}

func Test_ReadManifest(t *testing.T) {
	archive := writeMultiArchive(t, []string{"a", "b"}, map[string]map[string]string{
		"a": {"one.txt": "1"},
		"b": {"two.txt": "2"},
	})

	t.Run("manifest lists all volumes", func(t *testing.T) {
		manifest, err := dockergeneric.ReadManifest(bytes.NewReader(archive))
		require.NoError(t, err)
		require.Equal(t, []string{"a", "b"}, manifest.Volumes)
	})

	t.Run("archive without manifest returns ErrNoManifest", func(t *testing.T) {
		_, err := dockergeneric.ReadManifest(bytes.NewReader(buildTar(t, map[string]string{"x.txt": "y"})))
		require.ErrorIs(t, err, dockergeneric.ErrNoManifest)
	})

	t.Run("empty tar returns ErrNoManifest", func(t *testing.T) {
		_, err := dockergeneric.ReadManifest(bytes.NewReader(buildTar(t, nil)))
		require.ErrorIs(t, err, dockergeneric.ErrNoManifest)
	})

	t.Run("file named manifest.json in a subdirectory is not the manifest", func(t *testing.T) {
		_, err := dockergeneric.ReadManifest(bytes.NewReader(buildTar(t, map[string]string{
			"volumes/a/manifest.json": `{"version":1,"volumes":["a"]}`,
		})))
		require.ErrorIs(t, err, dockergeneric.ErrNoManifest)
	})

	t.Run("no tar archive returns an error", func(t *testing.T) {
		_, err := dockergeneric.ReadManifest(strings.NewReader(strings.Repeat("this is not a tar archive ", 100)))
		require.Error(t, err)
		require.NotErrorIs(t, err, dockergeneric.ErrNoManifest)
	})

	t.Run("nil reader returns an error", func(t *testing.T) {
		_, err := dockergeneric.ReadManifest(nil)
		require.Error(t, err)
	})
}

func Test_WriteMultiVolumeArchive(t *testing.T) {
	t.Run("all files are stored below volumes/<name>", func(t *testing.T) {
		archive := writeMultiArchive(t, []string{"a", "b"}, map[string]map[string]string{
			"a": {"one.txt": "1", "dir/two.txt": "2"},
			"b": {"three.txt": "3"},
		})

		files, _ := readTarEntries(t, archive)
		require.Equal(t, map[string]string{
			"manifest.json":         `{"version":1,"volumes":["a","b"]}`,
			"volumes/a/one.txt":     "1",
			"volumes/a/dir/two.txt": "2",
			"volumes/b/three.txt":   "3",
		}, files)
	})

	t.Run("empty volume is kept in the archive", func(t *testing.T) {
		archive := writeMultiArchive(t, []string{"empty"}, map[string]map[string]string{})

		_, dirs := readTarEntries(t, archive)
		require.Contains(t, dirs, "volumes/empty")

		manifest, err := dockergeneric.ReadManifest(bytes.NewReader(archive))
		require.NoError(t, err)
		require.Equal(t, []string{"empty"}, manifest.Volumes)
	})

	t.Run("long paths survive", func(t *testing.T) {
		longName := strings.Repeat("d", 60) + "/" + strings.Repeat("f", 60) + ".txt"
		archive := writeMultiArchive(t, []string{"a"}, map[string]map[string]string{
			"a": {longName: "long"},
		})

		files, _ := readTarEntries(t, archive)
		require.Equal(t, "long", files["volumes/a/"+longName])
	})

	t.Run("invalid volume name is rejected", func(t *testing.T) {
		var out bytes.Buffer
		err := dockergeneric.WriteMultiVolumeArchive(
			context.Background(),
			&out,
			[]string{"../etc"},
			func(string) (io.ReadCloser, error) {
				return io.NopCloser(bytes.NewReader(buildTar(t, nil))), nil
			},
		)
		require.Error(t, err)
	})

	t.Run("no volume names returns an error", func(t *testing.T) {
		var out bytes.Buffer
		err := dockergeneric.WriteMultiVolumeArchive(
			context.Background(),
			&out,
			[]string{},
			func(string) (io.ReadCloser, error) {
				return io.NopCloser(bytes.NewReader(buildTar(t, nil))), nil
			},
		)
		require.Error(t, err)
	})

	t.Run("error of the volume opener is returned", func(t *testing.T) {
		var out bytes.Buffer
		err := dockergeneric.WriteMultiVolumeArchive(
			context.Background(),
			&out,
			[]string{"a"},
			func(string) (io.ReadCloser, error) {
				return nil, errors.New("open failed")
			},
		)
		require.Error(t, err)
	})

	t.Run("nil writer returns an error", func(t *testing.T) {
		err := dockergeneric.WriteMultiVolumeArchive(
			context.Background(),
			nil,
			[]string{"a"},
			func(string) (io.ReadCloser, error) {
				return io.NopCloser(bytes.NewReader(buildTar(t, nil))), nil
			},
		)
		require.Error(t, err)
	})

	t.Run("nil opener returns an error", func(t *testing.T) {
		var out bytes.Buffer
		err := dockergeneric.WriteMultiVolumeArchive(context.Background(), &out, []string{"a"}, nil)
		require.Error(t, err)
	})
}

func Test_ExtractVolume(t *testing.T) {
	archive := writeMultiArchive(t, []string{"a", "ab", "b"}, map[string]map[string]string{
		"a":  {"one.txt": "1", "dir/two.txt": "2"},
		"ab": {"other.txt": "x"},
		"b":  {"three.txt": "3"},
	})

	t.Run("returns only the content of the requested volume", func(t *testing.T) {
		var out bytes.Buffer
		err := dockergeneric.ExtractVolume(context.Background(), bytes.NewReader(archive), "a", &out)
		require.NoError(t, err)

		files, _ := readTarEntries(t, out.Bytes())
		require.Equal(t, map[string]string{"one.txt": "1", "dir/two.txt": "2"}, files)
	})

	t.Run("volume with the same name prefix is not mixed up", func(t *testing.T) {
		var out bytes.Buffer
		err := dockergeneric.ExtractVolume(context.Background(), bytes.NewReader(archive), "ab", &out)
		require.NoError(t, err)

		files, _ := readTarEntries(t, out.Bytes())
		require.Equal(t, map[string]string{"other.txt": "x"}, files)
	})

	t.Run("manifest is not part of the extracted volume", func(t *testing.T) {
		var out bytes.Buffer
		err := dockergeneric.ExtractVolume(context.Background(), bytes.NewReader(archive), "b", &out)
		require.NoError(t, err)

		files, _ := readTarEntries(t, out.Bytes())
		require.NotContains(t, files, dockergeneric.ManifestFileName)
	})

	t.Run("volume not in archive results in an empty tar", func(t *testing.T) {
		var out bytes.Buffer
		err := dockergeneric.ExtractVolume(context.Background(), bytes.NewReader(archive), "missing", &out)
		require.NoError(t, err)

		files, dirs := readTarEntries(t, out.Bytes())
		require.Empty(t, files)
		require.Empty(t, dirs)
	})

	t.Run("hard link pointing outside of the volume is rejected", func(t *testing.T) {
		var malicious bytes.Buffer
		tw := tar.NewWriter(&malicious)
		require.NoError(t, tw.WriteHeader(&tar.Header{
			Typeflag: tar.TypeLink,
			Name:     "./volumes/a/link",
			Linkname: "./volumes/b/secret.txt",
			Mode:     0o644,
		}))
		require.NoError(t, tw.Close())

		var out bytes.Buffer
		err := dockergeneric.ExtractVolume(context.Background(), bytes.NewReader(malicious.Bytes()), "a", &out)
		require.Error(t, err)
	})

	t.Run("round trip of a long path", func(t *testing.T) {
		longName := strings.Repeat("d", 60) + "/" + strings.Repeat("f", 60) + ".txt"
		long := writeMultiArchive(t, []string{"a"}, map[string]map[string]string{
			"a": {longName: "long"},
		})

		var out bytes.Buffer
		err := dockergeneric.ExtractVolume(context.Background(), bytes.NewReader(long), "a", &out)
		require.NoError(t, err)

		files, _ := readTarEntries(t, out.Bytes())
		require.Equal(t, map[string]string{longName: "long"}, files)
	})

	t.Run("invalid volume name is rejected", func(t *testing.T) {
		var out bytes.Buffer
		err := dockergeneric.ExtractVolume(context.Background(), bytes.NewReader(archive), "../a", &out)
		require.Error(t, err)
	})

	t.Run("nil reader returns an error", func(t *testing.T) {
		var out bytes.Buffer
		err := dockergeneric.ExtractVolume(context.Background(), nil, "a", &out)
		require.Error(t, err)
	})

	t.Run("nil writer returns an error", func(t *testing.T) {
		err := dockergeneric.ExtractVolume(context.Background(), bytes.NewReader(archive), "a", nil)
		require.Error(t, err)
	})
}
