package dockergeneric

import (
	"archive/tar"
	"context"
	"encoding/json"
	"errors"
	"io"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

const (
	ManifestFileName = "manifest.json"
	ManifestVersion  = 1
)

// ErrNoManifest is returned when a multi volume archive does not tell which volumes it contains.
var ErrNoManifest = errors.New("archive does not contain a " + ManifestFileName)

var validVolumeName = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]*$`)

// Manifest is stored as manifest.json in the root of a multi volume archive.
type Manifest struct {
	Version int      `json:"version"`
	Volumes []string `json:"volumes"`
}

// CheckVolumeName validates a docker volume name. Names end up in archive paths
// and shell commands, so they are validated strictly.
func CheckVolumeName(volumeName string) error {
	if volumeName == "" {
		return tracederrors.TracedErrorEmptyString("volumeName")
	}

	if !validVolumeName.MatchString(volumeName) {
		return tracederrors.TracedErrorf("Invalid docker volume name '%s'.", volumeName)
	}

	return nil
}

// GetVolumeDirectory returns the directory of a volume inside a multi volume archive.
func GetVolumeDirectory(volumeName string) string {
	return "volumes/" + volumeName
}

// NormalizeEntryName returns the tar entry name without leading './' or '/' and without trailing '/'.
// The archive root results in an empty string.
func NormalizeEntryName(name string) string {
	return strings.TrimPrefix(path.Clean("/"+name), "/")
}

// MarshalManifest creates the content of the manifest.json.
func MarshalManifest(volumeNames []string) ([]byte, error) {
	ret, err := json.Marshal(Manifest{Version: ManifestVersion, Volumes: volumeNames})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Marshal manifest failed: %w", err)
	}

	return ret, nil
}

// ParseManifest parses and validates the content of a manifest.json.
func ParseManifest(content []byte) (*Manifest, error) {
	ret := &Manifest{}

	err := json.Unmarshal(content, ret)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Unmarshal manifest failed: %w", err)
	}

	if ret.Version != ManifestVersion {
		return nil, tracederrors.TracedErrorf("Unsupported manifest version '%d'.", ret.Version)
	}

	if len(ret.Volumes) == 0 {
		return nil, tracederrors.TracedErrorf("Manifest does not list any volume.")
	}

	for _, name := range ret.Volumes {
		err := CheckVolumeName(name)
		if err != nil {
			return nil, err
		}
	}

	return ret, nil
}

// ReadManifest scans the tar stream for the manifest.
// Returns ErrNoManifest (wrapped) if the archive has none.
func ReadManifest(r io.Reader) (*Manifest, error) {
	if r == nil {
		return nil, tracederrors.TracedErrorNil("r")
	}

	tr := tar.NewReader(r)
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			return nil, tracederrors.TracedErrorf("%w", ErrNoManifest)
		}
		if err != nil {
			return nil, tracederrors.TracedErrorf("Read tar archive failed: %w", err)
		}

		if NormalizeEntryName(hdr.Name) != ManifestFileName {
			continue
		}

		content, err := io.ReadAll(io.LimitReader(tr, 1<<20))
		if err != nil {
			return nil, tracederrors.TracedErrorf("Read manifest failed: %w", err)
		}

		return ParseManifest(content)
	}
}

// WriteMultiVolumeArchive writes the multi volume archive format to w.
// openVolumeTar must return the tar stream of the content of the given volume
// (volume content at the root of the stream).
func WriteMultiVolumeArchive(ctx context.Context, w io.Writer, volumeNames []string, openVolumeTar func(volumeName string) (io.ReadCloser, error)) error {
	if w == nil {
		return tracederrors.TracedErrorNil("w")
	}
	if openVolumeTar == nil {
		return tracederrors.TracedErrorNil("openVolumeTar")
	}
	if len(volumeNames) == 0 {
		return tracederrors.TracedErrorf("No volume names specified.")
	}

	logging.LogInfoByCtxf(ctx, "Write multi volume archive for '%d' volumes started.", len(volumeNames))

	manifest, err := MarshalManifest(volumeNames)
	if err != nil {
		return err
	}

	tw := tar.NewWriter(w)

	err = tw.WriteHeader(&tar.Header{
		Typeflag: tar.TypeReg,
		Name:     "./" + ManifestFileName,
		Mode:     0o644,
		Size:     int64(len(manifest)),
		ModTime:  time.Now(),
	})
	if err != nil {
		return tracederrors.TracedErrorf("Write manifest header failed: %w", err)
	}

	_, err = tw.Write(manifest)
	if err != nil {
		return tracederrors.TracedErrorf("Write manifest failed: %w", err)
	}

	for _, name := range volumeNames {
		err = writeVolume(tw, name, openVolumeTar)
		if err != nil {
			return err
		}
	}

	err = tw.Close()
	if err != nil {
		return tracederrors.TracedErrorf("Close tar writer failed: %w", err)
	}

	logging.LogInfoByCtxf(ctx, "Write multi volume archive for '%d' volumes finished.", len(volumeNames))

	return nil
}

func writeVolume(tw *tar.Writer, volumeName string, open func(string) (io.ReadCloser, error)) error {
	err := CheckVolumeName(volumeName)
	if err != nil {
		return err
	}

	rc, err := open(volumeName)
	if err != nil {
		return err
	}
	defer rc.Close()

	prefix := "./" + GetVolumeDirectory(volumeName)
	wroteRoot := false

	tr := tar.NewReader(rc)
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return tracederrors.TracedErrorf("Read tar of volume '%s' failed: %w", volumeName, err)
		}

		relative := NormalizeEntryName(hdr.Name)
		if relative == "" {
			if hdr.Typeflag != tar.TypeDir {
				continue
			}
			wroteRoot = true
			hdr.Name = prefix + "/"
		} else {
			hdr.Name = prefix + "/" + relative
			if hdr.Typeflag == tar.TypeDir {
				hdr.Name += "/"
			}
		}

		if hdr.Typeflag == tar.TypeLink {
			hdr.Linkname = prefix + "/" + NormalizeEntryName(hdr.Linkname)
		}

		// The new names can be longer than the original ones. Let the writer choose a suitable format.
		hdr.Format = tar.FormatUnknown

		err = tw.WriteHeader(hdr)
		if err != nil {
			return tracederrors.TracedErrorf("Write header '%s' failed: %w", hdr.Name, err)
		}

		_, err = io.Copy(tw, tr)
		if err != nil {
			return tracederrors.TracedErrorf("Write content of '%s' failed: %w", hdr.Name, err)
		}
	}

	if !wroteRoot { // Ensures empty volumes are present in the archive.
		err = tw.WriteHeader(&tar.Header{Typeflag: tar.TypeDir, Name: prefix + "/", Mode: 0o755, ModTime: time.Now()})
		if err != nil {
			return tracederrors.TracedErrorf("Write directory of volume '%s' failed: %w", volumeName, err)
		}
	}

	return nil
}

// ExtractVolume reads a multi volume archive from r and writes the content of
// the given volume as plain tar stream (volume content at the root) to w.
func ExtractVolume(ctx context.Context, r io.Reader, volumeName string, w io.Writer) error {
	if r == nil {
		return tracederrors.TracedErrorNil("r")
	}
	if w == nil {
		return tracederrors.TracedErrorNil("w")
	}

	err := CheckVolumeName(volumeName)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Extract volume '%s' from multi volume archive started.", volumeName)

	prefix := GetVolumeDirectory(volumeName)
	tr := tar.NewReader(r)
	tw := tar.NewWriter(w)

	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return tracederrors.TracedErrorf("Read tar archive failed: %w", err)
		}

		name := NormalizeEntryName(hdr.Name)
		switch {
		case name == prefix:
			hdr.Name = "./"
		case strings.HasPrefix(name, prefix+"/"):
			hdr.Name = strings.TrimPrefix(name, prefix+"/")
			if hdr.Typeflag == tar.TypeDir {
				hdr.Name += "/"
			}
		default:
			continue
		}

		if hdr.Typeflag == tar.TypeLink {
			link := NormalizeEntryName(hdr.Linkname)
			if !strings.HasPrefix(link, prefix+"/") {
				return tracederrors.TracedErrorf("Hard link '%s' points outside of volume '%s'.", hdr.Name, volumeName)
			}
			hdr.Linkname = strings.TrimPrefix(link, prefix+"/")
		}

		hdr.Format = tar.FormatUnknown

		err = tw.WriteHeader(hdr)
		if err != nil {
			return tracederrors.TracedErrorf("Write header '%s' failed: %w", hdr.Name, err)
		}

		_, err = io.Copy(tw, tr)
		if err != nil {
			return tracederrors.TracedErrorf("Write content of '%s' failed: %w", hdr.Name, err)
		}
	}

	err = tw.Close()
	if err != nil {
		return tracederrors.TracedErrorf("Close tar writer failed: %w", err)
	}

	logging.LogInfoByCtxf(ctx, "Extract volume '%s' from multi volume archive finished.", volumeName)

	return nil
}
