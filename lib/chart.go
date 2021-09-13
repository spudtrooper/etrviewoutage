package lib

func allPercentages(region string, seriesData regionSeriesData) (string, error) {
	t := `
	<html>
	<head>
		<title>All Percentages: {{.Region}}</title>
		<script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
		<script>
			google.charts.load('current', {packages: ['corechart', 'bar']});
			google.charts.setOnLoadCallback(drawMultSeries);
			
			function drawMultSeries() {
				var data = google.visualization.arrayToDataTable([
					['Parish', 'Percentage Affected'],
					{{range .Data}}
						['{{.County}}', {{.PercentAffected}}],
					{{end}}
				]);
			
				var options = {
					title: 'Population of Largest U.S. Cities',
					chartArea: {width: '50%', height: '100%'},
					hAxis: {
					title: 'All Percentages: {{.Region}}',
					minValue: 0
					},
					vAxis: {
					title: 'City'
					}
				};
			
				var chart = new google.visualization.BarChart(document.getElementById('chart_div'));
				chart.draw(data, options);
			}
		</script>
	</head>
	<body>
		<div id="chart_div"></div>
		<div>
			Last updated: <em>{{.LastUpdated}}</em>
		</div>
	</body>
`
	type regionData struct {
		County          string
		PercentAffected float64
	}
	var regionDatas []regionData
	for _, d := range seriesData {
		last := d.data[len(d.data)-1]
		percentAffected := float64(last.customersAffected) / float64(last.customersServed)
		regionDatas = append(regionDatas, regionData{
			County:          d.key.county,
			PercentAffected: percentAffected,
		})
	}
	data := struct {
		Region      string
		Data        []regionData
		LastUpdated string
	}{
		Region:      region,
		Data:        regionDatas,
		LastUpdated: lastUpdated(),
	}
	html, err := renderTemplate(t, "AllPercentages", data)
	if err != nil {
		return "", err
	}
	return html, nil
}

func percentageLineChart(series regionSeries) (string, error) {
	t := `
	<html>
	<head>
		<title>Percentages: {{.County}}</title>
		<script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
		<script>
			google.charts.load('current', {packages: ['corechart', 'line']});
			google.charts.setOnLoadCallback(drawBasic);
			
			function drawBasic() {		
			  {
				var data = new google.visualization.DataTable();
				data.addColumn('date', 'X');
				data.addColumn('number', 'Percentage Affected');
			
				data.addRows([
					{{range .Data}}
						[new Date({{.LastUpdateTime}}), {{.PercentAffected}}],
					{{end}}
				]);
			
				var options = {
					hAxis: {
						title: 'Time'
					},
					vAxis: {
						title: 'Percentage Affected'
					}
				};
			
				var chart = new google.visualization.LineChart(document.getElementById('chart_div'));
			
				chart.draw(data, options);
			  }
			  {
				var data = new google.visualization.DataTable();
				data.addColumn('date', 'X');
				data.addColumn('number', 'Affected');
				data.addColumn('number', 'Total');
			
				data.addRows([
					{{range .DataTotals}}
						[new Date({{.LastUpdateTime}}), {{.CustomersAffected}}, {{.CustomersServed}}],
					{{end}}
				]);
			
				var options = {
					hAxis: {
						title: 'Time'
					},
					vAxis: {
						title: 'Population'
					}
				};
			
				var chart = new google.visualization.LineChart(document.getElementById('chart_div_totals'));
			
				chart.draw(data, options);
			  }
			}
		</script>
	</head>
	<body>
		<div id="chart_div"></div>
		<div id="chart_div_totals"></div>
		<div>
			Last updated: <em>{{.LastUpdated}}</em>
		</div>
	</body>
`
	type regionData struct {
		LastUpdateTime  int
		PercentAffected float64
	}
	type regionDataTotal struct {
		LastUpdateTime    int
		CustomersAffected int
		CustomersServed   int
	}
	var regionDatas []regionData
	var regionDataTotals []regionDataTotal
	for _, d := range series.data {
		percentAffected := float64(d.customersAffected) / float64(d.customersServed)
		rd := regionData{
			LastUpdateTime:  d.lastUpdatedTime,
			PercentAffected: percentAffected,
		}
		regionDatas = append(regionDatas, rd)
		rdt := regionDataTotal{
			LastUpdateTime:    d.lastUpdatedTime,
			CustomersAffected: d.customersAffected,
			CustomersServed:   d.customersServed,
		}
		regionDataTotals = append(regionDataTotals, rdt)
	}
	data := struct {
		County      string
		Data        []regionData
		DataTotals  []regionDataTotal
		LastUpdated string
	}{
		County:      series.key.county,
		Data:        regionDatas,
		DataTotals:  regionDataTotals,
		LastUpdated: lastUpdated(),
	}
	html, err := renderTemplate(t, "Percentage", data)
	if err != nil {
		return "", err
	}
	return html, nil
}
