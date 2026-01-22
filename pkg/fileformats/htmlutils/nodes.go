package htmlutils

import (
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func GetNodeText(n *html.Node) string {
	if n == nil {
		return ""
	}

	var sb strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			sb.WriteString(strings.TrimSpace(c.Data))
		}
		if c.Type == html.CommentNode {
			rawContent := c.Data
			const prefix = "[CDATA["
			const suffix = "]]"
			if stringsutils.HasPrefixAndSuffix(rawContent, prefix, suffix) {
				sb.WriteString(stringsutils.TrimPrefixAndSuffix(rawContent, prefix, suffix))
			}
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

func GetAttributeValue(node *html.Node, namespace string, key string) (string, error) {
	if node == nil {
		return "", tracederrors.TracedErrorNil("node")
	}

	if key == "" {
		return "", tracederrors.TracedErrorEmptyString("key")
	}

	for _, attr := range node.Attr {
		if namespace != attr.Namespace {
			continue
		}

		if key != attr.Key {
			continue
		}

		return attr.Val, nil
	}

	return "", tracederrors.TracedErrorf("%w: Attribute '%s' in namespace '%s'.", ErrAttrNotFound, key, namespace)
}

func ParseFragmentAsHtmlNode(fragment string) (*html.Node, error) {
	context := &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Body,
		Data:     "body",
	}

	nodes, err := html.ParseFragment(strings.NewReader(fragment), context)
	if err != nil {
		panic(err)
	}

	return nodes[0], nil
}

func GetFirstChildNodeByTagName(node *html.Node, childNodeTagName string) (*html.Node, error) {
	if node == nil {
		return nil, tracederrors.TracedErrorNil("node")
	}

	if childNodeTagName == "" {
		return nil, tracederrors.TracedErrorEmptyString("childNodeTagName")
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type != html.ElementNode {
			continue
		}

		tag := strings.ToLower(c.Data)
		if tag == strings.ToLower(childNodeTagName) {
			return c, nil
		}
	}

	return nil, tracederrors.TracedErrorf("%w: child node tag='%s'.", ErrChildNodeNotFound, childNodeTagName)
}
