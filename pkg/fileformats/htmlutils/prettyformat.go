package htmlutils

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/fileformats/xmlutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"golang.org/x/net/html"
)

func prettyFormat(n *html.Node, indent string, builder *strings.Builder) {
	if n.Type == html.DocumentNode {
		if n.FirstChild != nil {
			prettyFormat(n.FirstChild, "", builder)
		}
	}
	if n.Type == html.TextNode {
		fmt.Fprintf(builder, "%s", n.Data)
	}
	if n.Type == html.ElementNode {
		if n.FirstChild == nil {
			// Empty element like <head></head>
			fmt.Fprintf(builder, "%s<%s></%s>\n", indent, n.Data, n.Data)
		} else {
			// Non empty element

			fmt.Fprintf(builder, "%s<%s", indent, n.Data)
			for _, attr := range n.Attr {
				fmt.Fprintf(builder, " %s=\"%s\"", attr.Key, attr.Val)
			}
			fmt.Fprintf(builder, ">")
			if n.FirstChild.Type != html.TextNode {
				fmt.Fprintf(builder, "\n")
			}

			prettyFormat(n.FirstChild, indent+"  ", builder)
			if n.FirstChild.Type == html.TextNode {
				fmt.Fprintf(builder, "</%s>\n", n.Data)
			} else {
				fmt.Fprintf(builder, "%s</%s>\n", indent, n.Data)
			}
		}

	}

	if n.NextSibling != nil {
		prettyFormat(n.NextSibling, indent, builder)
	}
}

func PrettyFormat(input string) (string, error) {
	doc, err := html.Parse(bytes.NewBufferString(input))
	if err != nil {
		return "", tracederrors.TracedErrorf("error parsing HTML: %w", err)
	}
	var builder strings.Builder

	err = html.Render(&builder, doc)
	if err != nil {
		return "", tracederrors.TracedErrorf("failed to render HTML back: %w", err)
	}

	formatted, err := xmlutils.PrettyFormat(builder.String())
	if err != nil {
		return "", err
	}

	if !strings.HasSuffix(formatted, "\n") {
		formatted += "\n"
	}

	return formatted, nil
}
