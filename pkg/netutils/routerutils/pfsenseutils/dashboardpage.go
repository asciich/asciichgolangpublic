package pfsenseutils

import (
	"context"
	"os"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/documentutils/htmldocument"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// The dashboard page refers to the main page shown in the pfSense UI entitled "Dashboard".

// Extracts the pfSense system name out of the dashboard HTML content.
func ParseDashboardPageAndGetSystemName(ctx context.Context, dashboardHtml []byte) (string, error) {
	if dashboardHtml == nil {
		return "", tracederrors.TracedErrorNil("dashboardHtml")
	}

	tables, err := htmldocument.ExtractTablesFromHTML(string(dashboardHtml))
	if err != nil {
		return "", err
	}

	var systemName string
	for _, table := range tables {
		nColumns, err := table.GetNumberOfColumns()
		if err != nil {
			return "", err
		}

		if nColumns <= 1 {
			continue
		}

		rows, err := table.GetRowsIncludingTitleRow()
		for _, row := range rows {
			entries, err := row.GetEntries()
			if err != nil {
				return "", err
			}

			if strings.EqualFold(entries[0], "Name") {
				systemName = strings.TrimSpace(entries[1])
				if systemName != "" {
					break
				}
			}
		}

		if systemName != "" {
			break
		}
	}

	if systemName == "" {
		os.WriteFile("againout.html", dashboardHtml, 0o644)
		return "", tracederrors.TracedError("Unable to extract system name from pfSense dashboard page. systemName is emtpy string after evaluation.")
	}

	logging.LogInfoByCtxf(ctx, "pfSense system name extracted from dashboard HTML content is '%s'.", systemName)

	return systemName, nil
}

func (r *Router) GetSystemName(ctx context.Context) (string, error) {
	url, err := r.GetUrl()
	if err != nil {
		return "", err
	}

	logging.LogInfoByCtxf(ctx, "Get system name of pfSense router '%s' started.", url)

	body, err := r.GetRequest(ctx, url)
	if err != nil {
		return "", err
	}

	systemName, err := ParseDashboardPageAndGetSystemName(ctx, body)
	if err != nil {
		return "", err
	}

	logging.LogInfoByCtxf(ctx, "The pfSense router '%s' is named '%s'.", url, systemName)

	logging.LogInfoByCtxf(ctx, "Get system name of pfSense router '%s' finished.", url)

	return systemName, err
}
