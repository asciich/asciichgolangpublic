package fstabconfigfile

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func ReadFromFile(ctx context.Context, path string) ([]*Entry, error) {
	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	content, err := nativefiles.ReadAsString(ctx, path, &filesoptions.ReadOptions{})
	if err != nil {
		return nil, err
	}

	return ParseFromString(content)
}

func ParseFromString(content string) ([]*Entry, error) {
	entries := []*Entry{}

	for _, line := range stringsutils.SplitLines(content, true) {
		line := stringsutils.RemoveComments(line)
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		line = stringsutils.RepeatReplaceAll(line, "\t", " ")
		line = stringsutils.RepeatReplaceAll(line, "  ", " ")

		splitted := strings.Split(line, " ")
		if len(splitted) != 6 {
			return nil, tracederrors.TracedErrorf("Unable to parse fstab line '%s'.", line)
		}

		entries = append(entries, &Entry{
			Device:  splitted[0],
			Dir:     splitted[1],
			Type:    splitted[2],
			Options: splitted[3],
			Dump:    splitted[4],
			Fsck:    splitted[5],
		})
	}

	return entries, nil
}
