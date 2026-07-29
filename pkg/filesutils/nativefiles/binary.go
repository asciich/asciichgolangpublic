package nativefiles

import (
	"context"
	"debug/elf"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// IsStaticallyLinkedBinary checks if the given path points to a statically linked binary.
// Returns true if the file is a statically linked binary, false otherwise.
// Returns an error if the file does not exist, is not a valid binary, or if the check fails.
func IsStaticallyLinkedBinary(ctx context.Context, path string) (bool, error) {
	err := ctx.Err()
	if err != nil {
		return false, err
	}

	if path == "" {
		return false, tracederrors.TracedErrorEmptyString("path")
	}

	if !Exists(contextutils.WithSilent(ctx), path) {
		return false, tracederrors.TracedErrorf("Path '%s' does not exist.", path)
	}

	if !IsFile(contextutils.WithSilent(ctx), path) {
		return false, tracederrors.TracedErrorf("Path '%s' is not a file.", path)
	}

	// Open the file as an ELF binary
	elfFile, err := elf.Open(path)
	if err != nil {
		// If the file is not an ELF binary, it's not a statically linked binary
		logging.LogInfoByCtxf(ctx, "'%s' is not a valid ELF binary: %v", path, err)
		return false, nil
	}
	defer elfFile.Close()

	// Check for PT_INTERP program header - indicates dynamic linking
	// Statically linked binaries don't have an interpreter
	hasInterp := false
	for _, prog := range elfFile.Progs {
		if prog.Type == elf.PT_INTERP {
			hasInterp = true
			break
		}
	}

	// Check for dynamic section - indicates dynamic linking
	// Statically linked binaries don't have a .dynamic section
	hasDynamic := false
	for _, section := range elfFile.Sections {
		if section.Type == elf.SHT_DYNAMIC {
			hasDynamic = true
			break
		}
	}

	// A binary is statically linked if it has no interpreter and no dynamic section
	isStaticallyLinked := !hasInterp && !hasDynamic

	if isStaticallyLinked {
		logging.LogInfoByCtxf(ctx, "'%s' is a statically linked binary.", path)
	} else {
		logging.LogInfoByCtxf(ctx, "'%s' is not a statically linked binary.", path)
	}

	return isStaticallyLinked, nil
}
