package lib

import (
	"path"
	"sort"
	"strings"
)

func outputCharts(dataDir string, dataMap map[string][]regionReportData) error {
	currentDir, err := mkdirAll(dataDir, "html", "current")
	if err != nil {
		return err
	}
	seriesData, err := readRegionReportDatas(dataMap)
	if err != nil {
		return err
	}

	for _, r := range regions {
		seriesDataFilteredByState := seriesData.filterByState(r)
		if len(seriesDataFilteredByState) == 0 {
			continue
		}
		sort.Slice(seriesDataFilteredByState, func(i, j int) bool {
			return strings.Compare(seriesDataFilteredByState[i].key.county, seriesDataFilteredByState[j].key.county) < 0
		})
		{
			html, err := allPercentages(r, seriesDataFilteredByState)
			if err != nil {
				return err
			}
			outFile := path.Join(currentDir, r+".html")
			if err := writeFile(outFile, []byte(html)); err != nil {
				return err
			}
		}
		stateHistoryDir, err := mkdirAll(dataDir, "html", "history", seriesDataFilteredByState[0].key.state)
		if err != nil {
			return err
		}
		{
			html, err := chartIndex(r, seriesDataFilteredByState)
			if err != nil {
				return err
			}
			outFile := path.Join(stateHistoryDir, "index.html")
			if err := writeFile(outFile, []byte(html)); err != nil {
				return err
			}
		}
		{
			html, err := statePercentAffectedTable(r, seriesDataFilteredByState)
			if err != nil {
				return err
			}
			outFile := path.Join(stateHistoryDir, "percent_affected_bydate.html")
			if err := writeFile(outFile, []byte(html)); err != nil {
				return err
			}
		}
		{
			html, err := statePercentDeltaTable(r, seriesDataFilteredByState)
			if err != nil {
				return err
			}
			outFile := path.Join(stateHistoryDir, "percent_delta_bydate.html")
			if err := writeFile(outFile, []byte(html)); err != nil {
				return err
			}
		}
	}
	for _, d := range seriesData {
		html, err := percentageLineChart(d)
		if err != nil {
			return err
		}
		historyDir, err := mkdirAll(dataDir, "html", "history", d.key.state)
		if err != nil {
			return err
		}
		outFile := path.Join(historyDir, d.key.county+".html")
		if err := writeFile(outFile, []byte(html)); err != nil {
			return err
		}
	}
	return nil
}
