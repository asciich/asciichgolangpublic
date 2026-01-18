package htmlutils

import (
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"golang.org/x/net/html"
)

func GetNodeText(n *html.Node) string {
	var sb strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			sb.WriteString(strings.TrimSpace(c.Data))
		}
	}
	return sb.String()
}

func FindBodyNode(node *html.Node) (*html.Node, error) {
	var body *html.Node
	var findBody func(*html.Node)
	findBody = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "body" {
			body = n
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findBody(c)
			if body != nil {
				break
			}
		}
	}
	findBody(node)

	if body == nil {
		return nil, tracederrors.TracedErrorf("%w", ErrNoHtmlBodyFound)
	}

	return body, nil
}
