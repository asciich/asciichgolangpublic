package commandexecutorfile_test

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
)

// createTestFileWithContent creates a file in a new temporary directory and returns its path.
func createTestFileWithContent(t *testing.T, fileName string, content string) string {
	t.Helper()

	filePath := filepath.Join(t.TempDir(), fileName)
	require.NoError(t, os.WriteFile(filePath, []byte(content), 0644))

	return filePath
}

// readTestFileContent returns the content of the given file as a string.
func readTestFileContent(t *testing.T, filePath string) string {
	t.Helper()

	content, err := os.ReadFile(filePath)
	require.NoError(t, err)

	return string(content)
}

func TestEndsWithLineBreak(t *testing.T) {
	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		ctx := getCtx()

		endsWithLineBreak, err := commandexecutorfile.EndsWithLineBreak(ctx, nil, "/tmp/some_file.txt")
		require.Error(t, err)
		require.False(t, endsWithLineBreak)
	})

	t.Run("empty path returns error", func(t *testing.T) {
		ctx := getCtx()

		endsWithLineBreak, err := commandexecutorfile.EndsWithLineBreak(ctx, commandexecutorexecoo.Exec(), "")
		require.Error(t, err)
		require.False(t, endsWithLineBreak)
	})

	t.Run("nonexistent file returns error", func(t *testing.T) {
		ctx := getCtx()

		endsWithLineBreak, err := commandexecutorfile.EndsWithLineBreak(ctx, commandexecutorexecoo.Exec(), "/tmp/this_file_does_not_exist_abc123xyz")
		require.Error(t, err)
		require.False(t, endsWithLineBreak)
	})

	t.Run("empty file returns false", func(t *testing.T) {
		ctx := getCtx()

		filePath := createTestFileWithContent(t, "empty.txt", "")

		endsWithLineBreak, err := commandexecutorfile.EndsWithLineBreak(ctx, commandexecutorexecoo.Exec(), filePath)
		require.NoError(t, err)
		require.False(t, endsWithLineBreak)
	})

	t.Run("file ending with line break returns true", func(t *testing.T) {
		ctx := getCtx()

		filePath := createTestFileWithContent(t, "with_newline.txt", "hello world\n")

		endsWithLineBreak, err := commandexecutorfile.EndsWithLineBreak(ctx, commandexecutorexecoo.Exec(), filePath)
		require.NoError(t, err)
		require.True(t, endsWithLineBreak)
	})

	t.Run("file without line break at the end returns false", func(t *testing.T) {
		ctx := getCtx()

		filePath := createTestFileWithContent(t, "without_newline.txt", "hello world")

		endsWithLineBreak, err := commandexecutorfile.EndsWithLineBreak(ctx, commandexecutorexecoo.Exec(), filePath)
		require.NoError(t, err)
		require.False(t, endsWithLineBreak)
	})

	t.Run("file with only a line break returns true", func(t *testing.T) {
		ctx := getCtx()

		filePath := createTestFileWithContent(t, "newline_only.txt", "\n")

		endsWithLineBreak, err := commandexecutorfile.EndsWithLineBreak(ctx, commandexecutorexecoo.Exec(), filePath)
		require.NoError(t, err)
		require.True(t, endsWithLineBreak)
	})

	t.Run("file with multiple lines ending with line break returns true", func(t *testing.T) {
		ctx := getCtx()

		filePath := createTestFileWithContent(t, "multi_line.txt", "line1\nline2\nline3\n")

		endsWithLineBreak, err := commandexecutorfile.EndsWithLineBreak(ctx, commandexecutorexecoo.Exec(), filePath)
		require.NoError(t, err)
		require.True(t, endsWithLineBreak)
	})

	t.Run("file with multiple lines without final line break returns false", func(t *testing.T) {
		ctx := getCtx()

		filePath := createTestFileWithContent(t, "multi_line_no_end.txt", "line1\nline2\nline3")

		endsWithLineBreak, err := commandexecutorfile.EndsWithLineBreak(ctx, commandexecutorexecoo.Exec(), filePath)
		require.NoError(t, err)
		require.False(t, endsWithLineBreak)
	})

	t.Run("file with windows line ending returns true", func(t *testing.T) {
		ctx := getCtx()

		filePath := createTestFileWithContent(t, "crlf.txt", "hello world\r\n")

		endsWithLineBreak, err := commandexecutorfile.EndsWithLineBreak(ctx, commandexecutorexecoo.Exec(), filePath)
		require.NoError(t, err)
		require.True(t, endsWithLineBreak)
	})

	t.Run("file with trailing space after line break returns false", func(t *testing.T) {
		ctx := getCtx()

		filePath := createTestFileWithContent(t, "trailing_space.txt", "hello world\n ")

		endsWithLineBreak, err := commandexecutorfile.EndsWithLineBreak(ctx, commandexecutorexecoo.Exec(), filePath)
		require.NoError(t, err)
		require.False(t, endsWithLineBreak)
	})

	t.Run("file ending with multibyte character returns false", func(t *testing.T) {
		ctx := getCtx()

		filePath := createTestFileWithContent(t, "umlaut.txt", "Grüezi mitenand ä")

		endsWithLineBreak, err := commandexecutorfile.EndsWithLineBreak(ctx, commandexecutorexecoo.Exec(), filePath)
		require.NoError(t, err)
		require.False(t, endsWithLineBreak)
	})

	t.Run("large file ending with line break returns true", func(t *testing.T) {
		ctx := getCtx()

		filePath := createTestFileWithContent(t, "large.txt", strings.Repeat("a", 1024*1024)+"\n")

		endsWithLineBreak, err := commandexecutorfile.EndsWithLineBreak(ctx, commandexecutorexecoo.Exec(), filePath)
		require.NoError(t, err)
		require.True(t, endsWithLineBreak)
	})

	t.Run("path with spaces ending with line break returns true", func(t *testing.T) {
		ctx := getCtx()

		filePath := createTestFileWithContent(t, "path with spaces.txt", "hello\n")

		endsWithLineBreak, err := commandexecutorfile.EndsWithLineBreak(ctx, commandexecutorexecoo.Exec(), filePath)
		require.NoError(t, err)
		require.True(t, endsWithLineBreak)
	})

	t.Run("/etc/passwd ends with line break", func(t *testing.T) {
		ctx := getCtx()

		endsWithLineBreak, err := commandexecutorfile.EndsWithLineBreak(ctx, commandexecutorexecoo.Exec(), "/etc/passwd")
		require.NoError(t, err)
		require.True(t, endsWithLineBreak)
	})
}

func TestAppendLine(t *testing.T) {
	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		ctx := getCtx()

		err := commandexecutorfile.AppendLine(ctx, nil, "/tmp/some_file.txt", "hello")
		require.Error(t, err)
	})

	t.Run("empty path returns error", func(t *testing.T) {
		ctx := getCtx()

		err := commandexecutorfile.AppendLine(ctx, commandexecutorexecoo.Exec(), "", "hello")
		require.Error(t, err)
	})

	t.Run("empty line returns error", func(t *testing.T) {
		ctx := getCtx()

		filePath := createTestFileWithContent(t, "file.txt", "existing\n")

		err := commandexecutorfile.AppendLine(ctx, commandexecutorexecoo.Exec(), filePath, "")
		require.Error(t, err)

		// The file must stay unchanged.
		require.Equal(t, "existing\n", readTestFileContent(t, filePath))
	})

	t.Run("nonexistent file is created", func(t *testing.T) {
		ctx := getCtx()

		filePath := filepath.Join(t.TempDir(), "new_file.txt")

		err := commandexecutorfile.AppendLine(ctx, commandexecutorexecoo.Exec(), filePath, "hello")
		require.NoError(t, err)
		require.Equal(t, "hello\n", readTestFileContent(t, filePath))
	})

	t.Run("append to empty file", func(t *testing.T) {
		ctx := getCtx()

		filePath := createTestFileWithContent(t, "empty.txt", "")

		err := commandexecutorfile.AppendLine(ctx, commandexecutorexecoo.Exec(), filePath, "hello")
		require.NoError(t, err)
		require.Equal(t, "hello\n", readTestFileContent(t, filePath))
	})

	t.Run("append to file ending with line break", func(t *testing.T) {
		ctx := getCtx()

		filePath := createTestFileWithContent(t, "with_newline.txt", "first\n")

		err := commandexecutorfile.AppendLine(ctx, commandexecutorexecoo.Exec(), filePath, "hello")
		require.NoError(t, err)
		require.Equal(t, "first\nhello\n", readTestFileContent(t, filePath))
	})

	t.Run("append to file without line break at the end", func(t *testing.T) {
		ctx := getCtx()

		filePath := createTestFileWithContent(t, "without_newline.txt", "first")

		err := commandexecutorfile.AppendLine(ctx, commandexecutorexecoo.Exec(), filePath, "hello")
		require.NoError(t, err)
		require.Equal(t, "first\nhello\n", readTestFileContent(t, filePath))
	})

	t.Run("line with line break at the end is written only with one line break", func(t *testing.T) {
		ctx := getCtx()

		filePath := createTestFileWithContent(t, "file.txt", "")

		err := commandexecutorfile.AppendLine(ctx, commandexecutorexecoo.Exec(), filePath, "hello\n")
		require.NoError(t, err)
		require.Equal(t, "hello\n", readTestFileContent(t, filePath))
	})

	t.Run("leading and trailing line breaks are removed from line", func(t *testing.T) {
		ctx := getCtx()

		filePath := createTestFileWithContent(t, "file.txt", "first\n")

		err := commandexecutorfile.AppendLine(ctx, commandexecutorexecoo.Exec(), filePath, "\n\nhello\n\n")
		require.NoError(t, err)
		require.Equal(t, "first\nhello\n", readTestFileContent(t, filePath))
	})

	t.Run("multiple appends add lines in order", func(t *testing.T) {
		ctx := getCtx()

		filePath := filepath.Join(t.TempDir(), "multi.txt")

		for _, line := range []string{"line1", "line2", "line3"} {
			err := commandexecutorfile.AppendLine(ctx, commandexecutorexecoo.Exec(), filePath, line)
			require.NoError(t, err)
		}

		require.Equal(t, "line1\nline2\nline3\n", readTestFileContent(t, filePath))
	})

	t.Run("appending the same line twice adds it twice", func(t *testing.T) {
		ctx := getCtx()

		filePath := createTestFileWithContent(t, "file.txt", "")

		require.NoError(t, commandexecutorfile.AppendLine(ctx, commandexecutorexecoo.Exec(), filePath, "hello"))
		require.NoError(t, commandexecutorfile.AppendLine(ctx, commandexecutorexecoo.Exec(), filePath, "hello"))

		require.Equal(t, "hello\nhello\n", readTestFileContent(t, filePath))
	})

	t.Run("line with spaces and special characters", func(t *testing.T) {
		ctx := getCtx()

		filePath := createTestFileWithContent(t, "file.txt", "")

		line := `export PATH="$HOME/bin:$PATH" # 'quoted' & ; | * ä`

		err := commandexecutorfile.AppendLine(ctx, commandexecutorexecoo.Exec(), filePath, line)
		require.NoError(t, err)
		require.Equal(t, line+"\n", readTestFileContent(t, filePath))
	})

	t.Run("path with spaces", func(t *testing.T) {
		ctx := getCtx()

		filePath := createTestFileWithContent(t, "path with spaces.txt", "first\n")

		err := commandexecutorfile.AppendLine(ctx, commandexecutorexecoo.Exec(), filePath, "hello")
		require.NoError(t, err)
		require.Equal(t, "first\nhello\n", readTestFileContent(t, filePath))
	})
}

func TestReadAsLines(t *testing.T) {
	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		ctx := getCtx()

		lines, err := commandexecutorfile.ReadAsLines(ctx, nil, "/tmp/some_file.txt")
		require.Error(t, err)
		require.Nil(t, lines)
	})

	t.Run("empty path returns error", func(t *testing.T) {
		ctx := getCtx()

		lines, err := commandexecutorfile.ReadAsLines(ctx, commandexecutorexecoo.Exec(), "")
		require.Error(t, err)
		require.Nil(t, lines)
	})

	t.Run("nonexistent file returns error", func(t *testing.T) {
		ctx := getCtx()

		lines, err := commandexecutorfile.ReadAsLines(ctx, commandexecutorexecoo.Exec(), "/tmp/this_file_does_not_exist_abc123xyz")
		require.Error(t, err)
		require.Nil(t, lines)
	})

	t.Run("single line without line break", func(t *testing.T) {
		ctx := getCtx()

		filePath := createTestFileWithContent(t, "single.txt", "hello")

		lines, err := commandexecutorfile.ReadAsLines(ctx, commandexecutorexecoo.Exec(), filePath)
		require.NoError(t, err)
		require.Equal(t, []string{"hello"}, lines)
	})

	t.Run("multiple lines without final line break", func(t *testing.T) {
		ctx := getCtx()

		filePath := createTestFileWithContent(t, "multi.txt", "line1\nline2\nline3")

		lines, err := commandexecutorfile.ReadAsLines(ctx, commandexecutorexecoo.Exec(), filePath)
		require.NoError(t, err)
		require.Equal(t, []string{"line1", "line2", "line3"}, lines)
	})

	t.Run("multiple lines with final line break", func(t *testing.T) {
		ctx := getCtx()

		filePath := createTestFileWithContent(t, "multi_newline.txt", "line1\nline2\nline3\n")

		lines, err := commandexecutorfile.ReadAsLines(ctx, commandexecutorexecoo.Exec(), filePath)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(lines), 3)
		require.Equal(t, []string{"line1", "line2", "line3"}, lines[:3])
	})

	t.Run("empty lines in the middle are kept", func(t *testing.T) {
		ctx := getCtx()

		filePath := createTestFileWithContent(t, "empty_lines.txt", "a\n\nb")

		lines, err := commandexecutorfile.ReadAsLines(ctx, commandexecutorexecoo.Exec(), filePath)
		require.NoError(t, err)
		require.Equal(t, []string{"a", "", "b"}, lines)
	})

	t.Run("lines with spaces are kept unchanged", func(t *testing.T) {
		ctx := getCtx()

		filePath := createTestFileWithContent(t, "spaces.txt", "hello world\nfoo bar baz")

		lines, err := commandexecutorfile.ReadAsLines(ctx, commandexecutorexecoo.Exec(), filePath)
		require.NoError(t, err)
		require.Equal(t, []string{"hello world", "foo bar baz"}, lines)
	})

	t.Run("lines written by AppendLine can be read back", func(t *testing.T) {
		ctx := getCtx()

		filePath := filepath.Join(t.TempDir(), "roundtrip.txt")

		for _, line := range []string{"line1", "line2"} {
			require.NoError(t, commandexecutorfile.AppendLine(ctx, commandexecutorexecoo.Exec(), filePath, line))
		}

		lines, err := commandexecutorfile.ReadAsLines(ctx, commandexecutorexecoo.Exec(), filePath)
		require.NoError(t, err)
		require.Contains(t, lines, "line1")
		require.Contains(t, lines, "line2")
	})

	t.Run("path with spaces", func(t *testing.T) {
		ctx := getCtx()

		filePath := createTestFileWithContent(t, "path with spaces.txt", "line1\nline2")

		lines, err := commandexecutorfile.ReadAsLines(ctx, commandexecutorexecoo.Exec(), filePath)
		require.NoError(t, err)
		require.Equal(t, []string{"line1", "line2"}, lines)
	})

	t.Run("/etc/passwd contains root user", func(t *testing.T) {
		ctx := getCtx()

		lines, err := commandexecutorfile.ReadAsLines(ctx, commandexecutorexecoo.Exec(), "/etc/passwd")
		require.NoError(t, err)
		require.True(t, slices.ContainsFunc(lines, func(line string) bool {
			return strings.HasPrefix(line, "root:")
		}))
	})
}

func TestEnsureLineInFile(t *testing.T) {
	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		ctx := getCtx()

		err := commandexecutorfile.EnsureLineInFile(ctx, nil, "/tmp/some_file.txt", "hello")
		require.Error(t, err)
	})

	t.Run("empty path returns error", func(t *testing.T) {
		ctx := getCtx()

		err := commandexecutorfile.EnsureLineInFile(ctx, commandexecutorexecoo.Exec(), "", "hello")
		require.Error(t, err)
	})

	t.Run("empty line returns error", func(t *testing.T) {
		ctx := getCtx()

		filePath := createTestFileWithContent(t, "file.txt", "existing\n")

		err := commandexecutorfile.EnsureLineInFile(ctx, commandexecutorexecoo.Exec(), filePath, "")
		require.Error(t, err)
		require.Equal(t, "existing\n", readTestFileContent(t, filePath))
	})

	t.Run("nonexistent file is created with line", func(t *testing.T) {
		ctx := getCtx()

		filePath := filepath.Join(t.TempDir(), "new_file.txt")

		err := commandexecutorfile.EnsureLineInFile(ctx, commandexecutorexecoo.Exec(), filePath, "hello")
		require.NoError(t, err)
		require.Equal(t, "hello\n", readTestFileContent(t, filePath))
	})

	t.Run("line is added to empty file", func(t *testing.T) {
		ctx := getCtx()

		filePath := createTestFileWithContent(t, "empty.txt", "")

		err := commandexecutorfile.EnsureLineInFile(ctx, commandexecutorexecoo.Exec(), filePath, "hello")
		require.NoError(t, err)
		require.Equal(t, "hello\n", readTestFileContent(t, filePath))
	})

	t.Run("missing line is appended", func(t *testing.T) {
		ctx := getCtx()

		filePath := createTestFileWithContent(t, "file.txt", "first\nsecond\n")

		err := commandexecutorfile.EnsureLineInFile(ctx, commandexecutorexecoo.Exec(), filePath, "hello")
		require.NoError(t, err)
		require.Equal(t, "first\nsecond\nhello\n", readTestFileContent(t, filePath))
	})

	t.Run("missing line is appended to file without final line break", func(t *testing.T) {
		ctx := getCtx()

		filePath := createTestFileWithContent(t, "file.txt", "first")

		err := commandexecutorfile.EnsureLineInFile(ctx, commandexecutorexecoo.Exec(), filePath, "hello")
		require.NoError(t, err)
		require.Equal(t, "first\nhello\n", readTestFileContent(t, filePath))
	})

	t.Run("existing line leaves file unchanged", func(t *testing.T) {
		ctx := getCtx()

		filePath := createTestFileWithContent(t, "file.txt", "first\nhello\nlast\n")

		err := commandexecutorfile.EnsureLineInFile(ctx, commandexecutorexecoo.Exec(), filePath, "hello")
		require.NoError(t, err)
		require.Equal(t, "first\nhello\nlast\n", readTestFileContent(t, filePath))
	})

	t.Run("existing last line without final line break leaves file unchanged", func(t *testing.T) {
		ctx := getCtx()

		filePath := createTestFileWithContent(t, "file.txt", "first\nhello")

		err := commandexecutorfile.EnsureLineInFile(ctx, commandexecutorexecoo.Exec(), filePath, "hello")
		require.NoError(t, err)
		require.Equal(t, "first\nhello", readTestFileContent(t, filePath))
	})

	t.Run("calling twice adds line only once", func(t *testing.T) {
		ctx := getCtx()

		filePath := createTestFileWithContent(t, "file.txt", "first\n")

		require.NoError(t, commandexecutorfile.EnsureLineInFile(ctx, commandexecutorexecoo.Exec(), filePath, "hello"))
		require.NoError(t, commandexecutorfile.EnsureLineInFile(ctx, commandexecutorexecoo.Exec(), filePath, "hello"))

		require.Equal(t, "first\nhello\n", readTestFileContent(t, filePath))
	})

	t.Run("partial match is not treated as present", func(t *testing.T) {
		ctx := getCtx()

		filePath := createTestFileWithContent(t, "file.txt", "hello world\n")

		err := commandexecutorfile.EnsureLineInFile(ctx, commandexecutorexecoo.Exec(), filePath, "hello")
		require.NoError(t, err)
		require.Equal(t, "hello world\nhello\n", readTestFileContent(t, filePath))
	})

	t.Run("line with leading whitespace is not equal to line without", func(t *testing.T) {
		ctx := getCtx()

		filePath := createTestFileWithContent(t, "file.txt", "  hello\n")

		err := commandexecutorfile.EnsureLineInFile(ctx, commandexecutorexecoo.Exec(), filePath, "hello")
		require.NoError(t, err)
		require.Equal(t, "  hello\nhello\n", readTestFileContent(t, filePath))
	})

	t.Run("multiple different lines are all ensured", func(t *testing.T) {
		ctx := getCtx()

		filePath := filepath.Join(t.TempDir(), "multi.txt")

		for _, line := range []string{"line1", "line2", "line1", "line3", "line2"} {
			require.NoError(t, commandexecutorfile.EnsureLineInFile(ctx, commandexecutorexecoo.Exec(), filePath, line))
		}

		require.Equal(t, "line1\nline2\nline3\n", readTestFileContent(t, filePath))
	})

	t.Run("path with spaces", func(t *testing.T) {
		ctx := getCtx()

		filePath := createTestFileWithContent(t, "path with spaces.txt", "first\n")

		require.NoError(t, commandexecutorfile.EnsureLineInFile(ctx, commandexecutorexecoo.Exec(), filePath, "hello"))
		require.NoError(t, commandexecutorfile.EnsureLineInFile(ctx, commandexecutorexecoo.Exec(), filePath, "hello"))

		require.Equal(t, "first\nhello\n", readTestFileContent(t, filePath))
	})
}
