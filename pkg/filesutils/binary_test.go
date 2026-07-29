package filesutils_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
)

func Test_IsStaticallyLinkedBinary(t *testing.T) {
	tests := []struct {
		name                 string
		path                 string
		setupFunc            func(t *testing.T, ctx context.Context) (path string, cleanup func())
		wantErr              bool
		wantStaticallyLinked bool
		skipOnSetupFailure   bool
	}{
		{
			name:                 "empty string",
			path:                 "",
			wantErr:              true,
			wantStaticallyLinked: false,
		},
		{
			name:                 "nonexisting file",
			path:                 "/this/file/does/not/exist",
			wantErr:              true,
			wantStaticallyLinked: false,
		},
		{
			name:                 "directory",
			path:                 "/etc",
			wantErr:              true,
			wantStaticallyLinked: false,
		},
		{
			name:                 "text file",
			path:                 "/etc/hosts",
			wantErr:              false,
			wantStaticallyLinked: false,
		},
		{
			name:                 "dynamically linked binary",
			path:                 "/bin/ls",
			wantErr:              false,
			wantStaticallyLinked: false,
		},
		{
			name: "statically linked binary",
			setupFunc: func(t *testing.T, ctx context.Context) (string, func()) {
				tempDir, err := tempfiles.CreateTempDir(ctx)
				require.NoError(t, err)

				binaryPath := filepath.Join(tempDir, "static_binary")
				sourcePath := filepath.Join(tempDir, "main.go")

				sourceCode := `package main
func main() {}
`
				err = nativefiles.WriteString(ctx, sourcePath, sourceCode)
				require.NoError(t, err)

				cmd := exec.Command("go", "build", "-ldflags", "-extldflags '-static'", "-o", binaryPath, sourcePath)
				cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
				_, err = cmd.CombinedOutput()
				require.NoError(t, err)

				cleanup := func() {
					nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})
				}

				return binaryPath, cleanup
			},
			wantErr:              false,
			wantStaticallyLinked: true,
			skipOnSetupFailure:   true,
		},
	}

	implementationNames := []string{
		"commandExecutorFileExec",
		"commandExecutorFileBash",
		"nativefilesoo",
	}

	for _, tt := range tests {
		for _, implementation := range implementationNames {
			t.Run(tt.name+"_"+implementation, func(t *testing.T) {
				ctx := getCtx()

				path := tt.path
				if tt.setupFunc != nil {
					var cleanup func()
					path, cleanup = tt.setupFunc(t, ctx)
					if cleanup != nil {
						defer cleanup()
					}
				}

				isStaticallyLinked, err := nativefiles.IsStaticallyLinkedBinary(ctx, path)

				if tt.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
				require.Equal(t, tt.wantStaticallyLinked, isStaticallyLinked)
			})
		}
	}
}

func Test_IsStaticallyLinkedBinary_WithContextCancellation(t *testing.T) {
	tests := []struct {
		name                 string
		path                 string
		ctxFunc              func() context.Context
		wantErr              bool
		wantStaticallyLinked bool
	}{
		{
			name: "cancelled context",
			path: "/bin/ls",
			ctxFunc: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			wantErr:              true,
			wantStaticallyLinked: false,
		},
	}

	implementationNames := []string{
		"commandExecutorFileExec",
		"commandExecutorFileBash",
		"nativefilesoo",
	}

	for _, tt := range tests {
		for _, implementation := range implementationNames {
			t.Run(tt.name+"_"+implementation, func(t *testing.T) {
				ctx := tt.ctxFunc()

				isStaticallyLinked, err := nativefiles.IsStaticallyLinkedBinary(ctx, tt.path)

				if tt.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
				require.Equal(t, tt.wantStaticallyLinked, isStaticallyLinked)
			})
		}
	}
}
