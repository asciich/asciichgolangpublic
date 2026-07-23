package htmldocument

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"

	"github.com/asciich/asciichgolangpublic/pkg/spreadsheet"
)

// ExtractTablesFromHTML parses the given HTML string and returns a slice of
// SpreadSheet objects, one for each <table> found in the document.
func ExtractTablesFromHTML(htmlContent string) ([]*spreadsheet.SpreadSheet, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var tables []*spreadsheet.SpreadSheet
	var parseErr error

	var findTables func(*html.Node)
	findTables = func(n *html.Node) {
		if parseErr != nil {
			return
		}
		if n.Type == html.ElementNode && n.Data == "table" {
			s, err := parseTable(n)
			if err != nil {
				parseErr = err
				return
			}
			tables = append(tables, s)
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findTables(c)
		}
	}

	findTables(doc)
	if parseErr != nil {
		return nil, parseErr
	}

	return tables, nil
}
func parseTable(tableNode *html.Node) (*spreadsheet.SpreadSheet, error) {
	s := spreadsheet.NewSpreadSheet()

	var rows []*html.Node
	collectRows(tableNode, &rows)

	titleSet := false
	for _, row := range rows {
		cells := extractCells(row)
		if !titleSet && rowContainsTH(row) {
			err := s.SetColumnTitles(cells)
			if err != nil {
				return nil, fmt.Errorf("failed to set column titles: %w", err)
			}
			titleSet = true
		} else {
			if !titleSet {
				// No header row found yet and this row has no <th> elements.
				// Initialize with empty column titles matching the row width.
				emptyTitles := make([]string, len(cells))
				err := s.SetColumnTitles(emptyTitles)
				if err != nil {
					return nil, fmt.Errorf("failed to set column titles: %w", err)
				}
				titleSet = true
			}
			err := s.AddRow(cells)
			if err != nil {
				return nil, fmt.Errorf("failed to add row: %w", err)
			}
		}
	}

	return s, nil
}

func collectRows(n *html.Node, rows *[]*html.Node) {
	if n.Type == html.ElementNode && n.Data == "tr" {
		*rows = append(*rows, n)
		return
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		collectRows(c, rows)
	}
}

func rowContainsTH(row *html.Node) bool {
	for c := row.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "th" {
			return true
		}
	}
	return false
}

func extractCells(row *html.Node) []string {
	var cells []string
	for c := row.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && (c.Data == "td" || c.Data == "th") {
			cells = append(cells, getTextContent(c))
		}
	}
	return cells
}

func getTextContent(n *html.Node) string {
	var buf strings.Builder
	var extract func(*html.Node)
	extract = func(node *html.Node) {
		if node.Type == html.TextNode {
			buf.WriteString(node.Data)
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}
	extract(n)
	return strings.TrimSpace(buf.String())
}
