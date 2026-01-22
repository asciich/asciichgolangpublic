package markdowndocument

import (
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/documentutils/basicdocument"
	"github.com/asciich/asciichgolangpublic/pkg/documentutils/documentinterfaces"
)

// ParseFromString parses a markdown string and populates the document
func ParseFromString(markdown string) (documentinterfaces.Document, error) {
	lines := strings.SplitSeq(markdown, "\n")
	var inVerbatim bool
	var verbatimContent []string

	ret := basicdocument.NewBasicDocument()

	for line := range lines {
		if line == "```" {
			if inVerbatim {
				// End of verbatim block
				if err := ret.AddVerbatimByString(strings.Join(verbatimContent, "\n")); err != nil {
					return nil, err
				}
				verbatimContent = nil
				inVerbatim = false
			} else {
				// Start of verbatim block
				inVerbatim = true
			}
			continue
		}

		if inVerbatim {
			verbatimContent = append(verbatimContent, line)
			continue
		}

		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "#### ") {
			if err := ret.AddSubSubSubTitleByString(strings.TrimPrefix(line, "#### ")); err != nil {
				return nil, err
			}
		} else if strings.HasPrefix(line, "### ") {
			if err := ret.AddSubSubTitleByString(strings.TrimPrefix(line, "### ")); err != nil {
				return nil, err
			}
		} else if strings.HasPrefix(line, "## ") {
			if err := ret.AddSubTitleByString(strings.TrimPrefix(line, "## ")); err != nil {
				return nil, err
			}
		} else if strings.HasPrefix(line, "# ") {
			if err := ret.AddTitleByString(strings.TrimPrefix(line, "# ")); err != nil {
				return nil, err
			}
		} else if !strings.HasPrefix(line, "|") {
			if err := ret.AddTextByString(line); err != nil {
				return nil, err
			}
		}
	}

	return ret, nil
}
