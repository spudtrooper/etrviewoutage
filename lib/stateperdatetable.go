package lib

import (
	"sort"
	"strings"
	"time"

	"github.com/go-errors/errors"
)

func statePercentAffectedTable(state string, data regionSeriesData) (string, error) {
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
					<h1>{{.State}} (% Affected)</h1>
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
											{{range .Dates}}
												<th data-sortable="true">{{.}}</th>
											{{end}}
										</tr>
									</thead>
									<tbody>
										{{range .Rows}}
											<tr> 
												<td scope="col" style="width: 250px">
													<a href="{{.Region}}.html">{{.Region}}</a>
												</td>
												{{range .PercentAffecteds}}
													<td scope="col">{{.}}</td>
												{{end}}												
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
	tData, err := makeStateSeriesTemplateData(state, data)
	if err != nil {
		return "", errors.Errorf("makeStateSeriesTemplateData: %v", err)
	}
	html, err := renderTemplate(t, "Percentage", tData)
	if err != nil {
		return "", errors.Errorf("renderTemplate(Percentage): %v", err)
	}
	return html, nil
}

func statePercentDeltaTable(state string, data regionSeriesData) (string, error) {
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
					<h1>{{.State}} (% Deltas)</h1>
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
											{{range .Dates}}
												<th data-sortable="true">{{.}}</th>
											{{end}}
										</tr>
									</thead>
									<tbody>
										{{range .Rows}}
											<tr> 
												<td scope="col" style="width: 250px">
													<a href="{{.Region}}.html">{{.Region}}</a>
												</td>
												{{range .PercentDeltas}}
													<td scope="col">{{.}}</td>
												{{end}}												
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
	tData, err := makeStateSeriesTemplateData(state, data)
	if err != nil {
		return "", errors.Errorf("makeStateSeriesTemplateData: %v", err)
	}
	html, err := renderTemplate(t, "PercentDelta", tData)
	if err != nil {
		return "", errors.Errorf("renderTemplate(PercentDelta): %v", err)
	}
	return html, nil
}

type stateSeriesTemplateRow struct {
	Region           string
	PercentAffecteds []string
	PercentDeltas    []string
}

type stateSeriesTemplateData struct {
	State       string
	Rows        []stateSeriesTemplateRow
	LastUpdated string
	Dates       []string
}

func makeStateSeriesTemplateData(state string, data regionSeriesData) (stateSeriesTemplateData, error) {
	dateFromLastUpdatedTime := func(d regionData) string {
		return time.Unix(int64(d.lastUpdatedTime/1000), 0).Format("1/02")
	}
	dateMap := map[string]bool{}
	for _, row := range data {
		for _, d := range row.data {
			dateMap[dateFromLastUpdatedTime(d)] = true
		}
	}
	var dates []string
	for d := range dateMap {
		dates = append(dates, d)
	}
	sort.Strings(dates)

	var rows []stateSeriesTemplateRow
	for _, d := range data {
		regionDatasByDate := map[string][]regionData{}
		for _, d := range d.data {
			date := dateFromLastUpdatedTime(d)
			regionDatas, ok := regionDatasByDate[date]
			if !ok {
				regionDatas = []regionData{}
			}
			regionDatas = append(regionDatas, d)
			regionDatasByDate[date] = regionDatas
		}

		lastRegionDataByDate := map[string]regionData{}
		for d, ds := range regionDatasByDate {
			sort.Slice(ds, func(i, j int) bool {
				return ds[i].lastUpdatedTime < ds[j].lastUpdatedTime
			})
			last := ds[len(ds)-1]
			lastRegionDataByDate[d] = last
		}

		var percentAffecteds []string
		var percentDeltas []string
		var lastPercentage float64
		for _, date := range dates {
			r, ok := lastRegionDataByDate[date]
			var percentAffected float64
			if !ok {
				percentAffected = lastPercentage
			} else {
				percentAffected = r.percentAffected()
			}
			percentAffecteds = append(percentAffecteds, formatPercentage(percentAffected))
			percentDelta := "--"
			if lastPercentage != 0 {
				percentDelta = formatPercentage((percentAffected - lastPercentage) / lastPercentage)
			}
			percentDeltas = append(percentDeltas, percentDelta)

			lastPercentage = percentAffected
		}
		r := stateSeriesTemplateRow{
			Region:           d.key.county,
			PercentAffecteds: percentAffecteds,
			PercentDeltas:    percentDeltas,
		}
		rows = append(rows, r)
	}
	sort.Slice(rows, func(i, j int) bool {
		return strings.Compare(rows[i].Region, rows[j].Region) < 0
	})

	tData := stateSeriesTemplateData{
		State:       state,
		Rows:        rows,
		LastUpdated: lastUpdated(),
		Dates:       dates,
	}
	return tData, nil
}
