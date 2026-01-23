package htmlutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/fileformats/htmlutils"
	"golang.org/x/net/html"
)

func Test_GetNodeText(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		text := htmlutils.GetNodeText(nil)
		require.Empty(t, text)
	})

	t.Run("paragraph", func(t *testing.T) {
		node, err := htmlutils.ParseFragmentAsHtmlNode("<p>this is the text</p>")
		require.NoError(t, err)

		text := htmlutils.GetNodeText(node)
		require.EqualValues(t, "this is the text", text)
	})

	t.Run("ac:plain-text", func(t *testing.T) {
		node, err := htmlutils.ParseFragmentAsHtmlNode("<ac:plain-text-body><!--[CDATA[hello world]]--></ac:plain-text-body>")
		require.NoError(t, err)

		text := htmlutils.GetNodeText(node)
		require.EqualValues(t, "hello world", text)
	})
}

func Test_GetNodeAttributeValue(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		node := &html.Node{}
		value, err := htmlutils.GetAttributeValue(node, "", "key")
		require.Error(t, err)
		require.True(t, htmlutils.IsErrAttrNotFound(err))
		require.Empty(t, value)
	})

	t.Run("no namespace", func(t *testing.T) {
		node := &html.Node{
			Attr: []html.Attribute{
				{
					Namespace: "",
					Key:       "key",
					Val:       "value",
				},
			},
		}
		value, err := htmlutils.GetAttributeValue(node, "", "key")
		require.NoError(t, err)
		require.EqualValues(t, "value", value)
	})

	t.Run("with namespace", func(t *testing.T) {
		node := &html.Node{
			Attr: []html.Attribute{
				{
					Namespace: "ns",
					Key:       "key",
					Val:       "value",
				},
			},
		}
		value, err := htmlutils.GetAttributeValue(node, "ns", "key")
		require.NoError(t, err)
		require.EqualValues(t, "value", value)
	})

	t.Run("with namespace2", func(t *testing.T) {
		node := &html.Node{
			Attr: []html.Attribute{
				{
					Namespace: "",
					Key:       "key",
					Val:       "value",
				},
				{
					Namespace: "ns",
					Key:       "key",
					Val:       "value1",
				},
				{
					Namespace: "namespace",
					Key:       "key",
					Val:       "value2",
				},
				{
					Namespace: "namespace",
					Key:       "anotherkey",
					Val:       "value3",
				},
			},
		}
		value, err := htmlutils.GetAttributeValue(node, "namespace", "key")
		require.NoError(t, err)
		require.EqualValues(t, "value2", value)
	})
}

func Test_GetChildNodeByTagName(t *testing.T) {
	t.Run("nil, empty", func(t *testing.T) {
		node, err := htmlutils.GetFirstChildNodeByTagName(nil, "")
		require.Error(t, err)
		require.Nil(t, node)
	})

	t.Run("empty", func(t *testing.T) {
		node, err := htmlutils.GetFirstChildNodeByTagName(&html.Node{}, "")
		require.Error(t, err)
		require.Nil(t, node)
	})

	t.Run("not found", func(t *testing.T) {
		node, err := htmlutils.GetFirstChildNodeByTagName(&html.Node{}, "unknown")
		require.Error(t, err)
		require.True(t, htmlutils.IsErrChildNodeNotFound(err))
		require.Nil(t, node)
	})

	t.Run("atlassian confluence wiki code block", func(t *testing.T) {
		rawHtml := `<ac:structured-macro ac:name="code" ac:schema-version="1" ac:macro-id="12345678-abcd-1234-abcd-12345678911">
    <ac:parameter ac:name="language">shell</ac:parameter>
    <ac:plain-text-body>
    <!--[CDATA[echo "hello world!"]]--></ac:plain-text-body>
</ac:structured-macro>
	`
		node, err := htmlutils.ParseFragmentAsHtmlNode(rawHtml)
		require.NoError(t, err)

		child, err := htmlutils.GetFirstChildNodeByTagName(node, "ac:parameter")
		require.NoError(t, err)
		require.NotNil(t, child)
		require.EqualValues(t, "ac:parameter", child.Data)

		child, err = htmlutils.GetFirstChildNodeByTagName(node, "ac:plain-text-body")
		require.NoError(t, err)
		require.NotNil(t, child)
		require.EqualValues(t, "ac:plain-text-body", child.Data)
		require.EqualValues(t, "echo \"hello world!\"", htmlutils.GetNodeText(child))
	})
}
