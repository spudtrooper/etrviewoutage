package lib

import (
	"sort"
	"strings"
)

func chartIndex(state string, data regionSeriesData) (string, error) {
	t := `
		<html>
			<head>
				<link rel="stylesheet" type="text/css" href="https://getbootstrap.com/docs/4.0/dist/css/bootstrap.min.css">
				<link rel="stylesheet" type="text/css" href="https://unpkg.com/bootstrap-table@1.18.2/dist/bootstrap-table.min.css">
				<script src="https://cdnjs.cloudflare.com/ajax/libs/jquery/3.5.1/jquery.min.js" integrity="sha512-bLT0Qm9VnAYZDflyKcBaQ2gg0hSYNQrJ8RilYldYQ1FxQYoCLtUjuuRuZo+fjqhx/qtq/1itJ0C2ejDxltZVFg==" crossorigin="anonymous"></script>
				<script src="https://getbootstrap.com/docs/4.0/dist/js/bootstrap.min.js"></script>			
				<script src="https://unpkg.com/bootstrap-table@1.18.2/dist/bootstrap-table.min.js"></script>
				<style>
					.headhead {
						background: #003399;
						color: #fff;
					}				
				</style>
				<script>
					$(document).ready(function () {
						$('.sortable').bootstrapTable('refreshOptions', {
							sortable: true,
						});
					});
				</script>				
			</head>
			<body>
				<div class="container-fluid">			
					<h1>{{.State}}</h1>
					<div>
						<a href="../../../">home</a>
					</div>
					<div style="margin-bottom: 10px">
						Last updated: <em>{{.LastUpdated}}</em>
					</div>
					<table border=0 style="border: none; width:100%; height:100%">
						<tr>
							<td style="vertical-align:top; width:600px; overflow:auto">
								<table data-toggle="table" class="sortable table table-striped table-bordered table-sm">
									<thead>
										<tr class="headhead">
											<th data-sortable="true">Parish</th>
											<th data-sortable="true">Total</th>
											<th data-sortable="true">Affected</th>
											<th data-sortable="true">% Affected</th>
										</tr>
									</thead>
									<tbody>
										{{range .Rows}}
											<tr> 
												<td scope="col" style="width: 250px">
													<a href="javascript:document.getElementById('main').src='{{.Region}}.html';void(0);">{{.Region}}</a>
													<a href="{{.Region}}.html" target="_"><img style="width:20px; height:20px" src="data:image/svg+xml;base64,PHN2ZyBoZWlnaHQ9JzMwMHB4JyB3aWR0aD0nMzAwcHgnICBmaWxsPSIjMDAwMDAwIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIGRhdGEtbmFtZT0iTGF5ZXIgMSIgdmlld0JveD0iMCAwIDEwMCAxMDAiIHg9IjBweCIgeT0iMHB4Ij48dGl0bGU+MTc8L3RpdGxlPjxnIGRhdGEtbmFtZT0iR3JvdXAiPjxwYXRoIGRhdGEtbmFtZT0iUGF0aCIgZD0iTTgyLjEsMTMuMmgtMjF2OS4ySDcxTDQwLjEsNTMuNGw2LjUsNi41TDc3LjUsMjl2OS44aDkuMnYtMjFBNC42LDQuNiwwLDAsMCw4Mi4xLDEzLjJaIj48L3BhdGg+PHBhdGggZGF0YS1uYW1lPSJQYXRoIiBkPSJNNjguNiw3Ny41SDIyLjVWMzEuNEg1NC4xVjIyLjFIMTcuOWE0LjYsNC42LDAsMCwwLTQuNiw0LjZWODIuMWE0LjYsNC42LDAsMCwwLDQuNiw0LjZINzMuMmE0LjYsNC42LDAsMCwwLDQuNi00LjZWNDUuOUg2OC42WiI+PC9wYXRoPjwvZz48L3N2Zz4="></a>
												</td>
												<td scope="col">{{.CustomersServiced}}</td>
												<td scope="col">{{.CustomersAffected}}</td>
												<td scope="col">{{.PercentAffected}}</td>
											</tr>
										{{end}}
									</tbody>
								</table>
							</td>
							<td style="vertical-align:top; overflow:auto">
								<iframe id="main" src="" style="border: none; width:100%; height:100%"></iframe>
							</td>
						</tr>
					</table>
				</div>
			</body>
		</html>
		`
	type row struct {
		Region            string
		CustomersServiced int
		CustomersAffected int
		PercentAffected   string
	}
	var rows []row
	for _, d := range data {
		ds := []regionData{}
		ds = append(ds, d.data...)
		sort.Slice(ds, func(i, j int) bool {
			return ds[i].lastUpdatedTime < ds[j].lastUpdatedTime
		})
		last := ds[len(ds)-1]
		r := row{
			Region:            d.key.county,
			CustomersServiced: last.customersServed,
			CustomersAffected: last.customersAffected,
			PercentAffected:   formatPercentage(last.percentAffected()),
		}
		rows = append(rows, r)
	}
	sort.Slice(rows, func(i, j int) bool {
		return strings.Compare(rows[i].Region, rows[j].Region) < 0
	})
	tData := struct {
		State       string
		Rows        []row
		LastUpdated string
	}{
		State:       state,
		Rows:        rows,
		LastUpdated: lastUpdated(),
	}
	html, err := renderTemplate(t, "Percentage", tData)
	if err != nil {
		return "", err
	}
	return html, nil
}
