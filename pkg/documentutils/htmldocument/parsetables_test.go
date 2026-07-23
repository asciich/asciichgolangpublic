package htmldocument

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractTablesFromHTML_SingleTable(t *testing.T) {
	htmlContent := `
<!DOCTYPE html>
<html>
<body>
<table>
  <tr>
    <th>Name</th>
    <th>Age</th>
    <th>City</th>
  </tr>
  <tr>
    <td>Alice</td>
    <td>30</td>
    <td>Zurich</td>
  </tr>
  <tr>
    <td>Bob</td>
    <td>25</td>
    <td>Bern</td>
  </tr>
</table>
</body>
</html>
`

	tables, err := ExtractTablesFromHTML(htmlContent)
	require.NoError(t, err)
	require.EqualValues(t, 1, len(tables))

	spreadSheet := tables[0]

	numberOfColumns, err := spreadSheet.GetNumberOfColumns()
	require.NoError(t, err)
	require.EqualValues(t, 3, numberOfColumns)

	numberOfRows, err := spreadSheet.GetNumberOfRows()
	require.NoError(t, err)
	require.EqualValues(t, 2, numberOfRows)

	columnTitle0, err := spreadSheet.GetColumnTitleAtIndexAsString(0)
	require.NoError(t, err)
	require.EqualValues(t, "Name", columnTitle0)

	columnTitle1, err := spreadSheet.GetColumnTitleAtIndexAsString(1)
	require.NoError(t, err)
	require.EqualValues(t, "Age", columnTitle1)

	columnTitle2, err := spreadSheet.GetColumnTitleAtIndexAsString(2)
	require.NoError(t, err)
	require.EqualValues(t, "City", columnTitle2)

	cell00, err := spreadSheet.GetCellValueAsString(0, 0)
	require.NoError(t, err)
	require.EqualValues(t, "Alice", cell00)

	cell01, err := spreadSheet.GetCellValueAsString(0, 1)
	require.NoError(t, err)
	require.EqualValues(t, "30", cell01)

	cell02, err := spreadSheet.GetCellValueAsString(0, 2)
	require.NoError(t, err)
	require.EqualValues(t, "Zurich", cell02)

	cell10, err := spreadSheet.GetCellValueAsString(1, 0)
	require.NoError(t, err)
	require.EqualValues(t, "Bob", cell10)

	cell11, err := spreadSheet.GetCellValueAsString(1, 1)
	require.NoError(t, err)
	require.EqualValues(t, "25", cell11)

	cell12, err := spreadSheet.GetCellValueAsString(1, 2)
	require.NoError(t, err)
	require.EqualValues(t, "Bern", cell12)
}

func TestExtractTablesFromHTML_KeyValueTable(t *testing.T) {
	htmlContent := `
<html>
<body>
<table class="table table-hover table-striped table-condensed">
	<tbody>
		<tr>
			<th>Name</th>
			<td>router.example.com</td>
		</tr>
		<tr>
			<th>User</th>
			<td>admin@192.168.1.123 (Local Database)</td>
		</tr>
	</tbody>
</table>
</body>
</html>
`

	tables, err := ExtractTablesFromHTML(htmlContent)
	require.NoError(t, err)
	require.EqualValues(t, 1, len(tables))

	spreadSheet := tables[0]

	numberOfColumns, err := spreadSheet.GetNumberOfColumns()
	require.NoError(t, err)
	require.EqualValues(t, 2, numberOfColumns)

	numberOfRows, err := spreadSheet.GetNumberOfRows()
	require.NoError(t, err)
	require.EqualValues(t, 1, numberOfRows)

	columnTitle0, err := spreadSheet.GetColumnTitleAtIndexAsString(0)
	require.NoError(t, err)
	require.EqualValues(t, "Name", columnTitle0)

	columnTitle1, err := spreadSheet.GetColumnTitleAtIndexAsString(1)
	require.NoError(t, err)
	require.EqualValues(t, "router.example.com", columnTitle1)

	cell00, err := spreadSheet.GetCellValueAsString(0, 0)
	require.NoError(t, err)
	require.EqualValues(t, "User", cell00)

	cell01, err := spreadSheet.GetCellValueAsString(0, 1)
	require.NoError(t, err)
	require.EqualValues(t, "admin@192.168.1.123 (Local Database)", cell01)
}

func TestExtractTablesFromHTML_NoHeaderRow(t *testing.T) {
	htmlContent := `
<html>
<body>
<div class="col-md-6" id="widgets-col2">
	<div class="panel panel-default" id="widget-interfaces-0">
		<div class="panel-heading">
			<h2 class="panel-title"><a href="status_interfaces.php">Interfaces</a></h2>
		</div>
		<div id="widget-interfaces-0_panel-body" class="panel-body collapse in">
			<div class="table-responsive" id="ifaces_status_interfaces-0">
				<table class="table table-striped table-hover table-condensed">
					<tbody>
						<tr>
							<td><a href="/interfaces.php?if=wan">WAN</a></td>
							<td>up</td>
							<td>1000baseT &lt;full-duplex&gt;</td>
							<td>192.168.1.101</td>
						</tr>
						<tr>
							<td><a href="/interfaces.php?if=lan">LAN</a></td>
							<td>up</td>
							<td>1000baseT &lt;full-duplex&gt;</td>
							<td>192.168.10.1</td>
						</tr>
						<tr>
							<td><a href="/interfaces.php?if=opt1">OPT</a></td>
							<td>no carrier</td>
							<td>none</td>
							<td>n/a</td>
						</tr>
					</tbody>
				</table>
			</div>
		</div>
	</div>
</div>
</body>
</html>
`

	tables, err := ExtractTablesFromHTML(htmlContent)
	require.NoError(t, err)
	require.EqualValues(t, 1, len(tables))

	spreadSheet := tables[0]

	numberOfColumns, err := spreadSheet.GetNumberOfColumns()
	require.NoError(t, err)
	require.EqualValues(t, 4, numberOfColumns)

	numberOfRows, err := spreadSheet.GetNumberOfRows()
	require.NoError(t, err)
	require.EqualValues(t, 3, numberOfRows)

	cell00, err := spreadSheet.GetCellValueAsString(0, 0)
	require.NoError(t, err)
	require.EqualValues(t, "WAN", cell00)

	cell01, err := spreadSheet.GetCellValueAsString(0, 1)
	require.NoError(t, err)
	require.EqualValues(t, "up", cell01)

	cell02, err := spreadSheet.GetCellValueAsString(0, 2)
	require.NoError(t, err)
	require.EqualValues(t, "1000baseT <full-duplex>", cell02)

	cell03, err := spreadSheet.GetCellValueAsString(0, 3)
	require.NoError(t, err)
	require.EqualValues(t, "192.168.1.101", cell03)

	cell10, err := spreadSheet.GetCellValueAsString(1, 0)
	require.NoError(t, err)
	require.EqualValues(t, "LAN", cell10)

	cell11, err := spreadSheet.GetCellValueAsString(1, 1)
	require.NoError(t, err)
	require.EqualValues(t, "up", cell11)

	cell12, err := spreadSheet.GetCellValueAsString(1, 2)
	require.NoError(t, err)
	require.EqualValues(t, "1000baseT <full-duplex>", cell12)

	cell13, err := spreadSheet.GetCellValueAsString(1, 3)
	require.NoError(t, err)
	require.EqualValues(t, "192.168.10.1", cell13)

	cell20, err := spreadSheet.GetCellValueAsString(2, 0)
	require.NoError(t, err)
	require.EqualValues(t, "OPT", cell20)

	cell21, err := spreadSheet.GetCellValueAsString(2, 1)
	require.NoError(t, err)
	require.EqualValues(t, "no carrier", cell21)

	cell22, err := spreadSheet.GetCellValueAsString(2, 2)
	require.NoError(t, err)
	require.EqualValues(t, "none", cell22)

	cell23, err := spreadSheet.GetCellValueAsString(2, 3)
	require.NoError(t, err)
	require.EqualValues(t, "n/a", cell23)
}
