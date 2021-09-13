package lib

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-errors/errors"
)

// {"state":"L","county":"ASCENSION","customersServed":43240,"customersAffected":34490,"lastUpdatedTime":1630424672300,"latitude":30.258093385214007,"longitude":-90.84015564202333}
type regionRawData struct {
	State             string  `json:"state"`
	County            string  `json:"county"`
	CustomersServed   int     `json:"customersServed"`
	CustomersAffected int     `json:"customersAffected"`
	LastUpdatedTime   int     `json:"lastUpdatedTime"`
	Latitude          float64 `json:"latitude"`
	Longitude         float64 `json:"longitude"`
	sourceFile        string
}

var (
	stateNames = map[string]string{
		"L": "Louisiana",
		"M": "Mississippi",
		"T": "Texas",
		"A": "Alabama",
	}
)

func (r *regionRawData) toMetadata() regionMetadata {
	return regionMetadata{
		state:      r.State,
		stateNamme: stateNames[r.State],
		county:     r.County,
		latitude:   r.Latitude,
		longitude:  r.Longitude,
	}
}

func (r *regionRawData) toData() regionData {
	return regionData{
		lastUpdatedTime:   r.LastUpdatedTime,
		customersServed:   r.CustomersServed,
		customersAffected: r.CustomersAffected,
	}
}

type regionMetadata struct {
	state      string
	stateNamme string
	county     string
	latitude   float64
	longitude  float64
}

type regionData struct {
	lastUpdatedTime   int
	customersServed   int
	customersAffected int
}

func (r *regionData) percentAffected() float64 {
	return float64(r.customersAffected) / float64(r.customersServed)
}

type regionReportData struct {
	date       time.Time
	key        regionMetadata
	data       regionData
	sourceFile string
}

type regionSeries struct {
	key  regionMetadata
	data []regionData
}

type regionSeriesData []regionSeries

func (r *regionSeriesData) filterByState(state string) regionSeriesData {
	var res regionSeriesData
	for _, d := range *r {
		if strings.HasPrefix(state, d.key.state) {
			res = append(res, d)
		}
	}
	return res
}

func readRegionReportDatas(dataMap map[string][]regionReportData) (regionSeriesData, error) {
	var res regionSeriesData
	for _, datas := range dataMap {
		key := datas[0].key
		var regionDatas []regionData
		for _, d := range datas {
			regionDatas = append(regionDatas, d.data)
		}
		series := regionSeries{
			key:  key,
			data: regionDatas,
		}
		res = append(res, series)
	}
	return res, nil
}

func readRegionDataMap(dataDir string) (map[string][]regionReportData, []time.Time, error) {
	regionsDir := path.Join(dataDir, "regions")
	regsionsPattern := path.Join(regionsDir, "*")
	regionDirs, err := filepath.Glob(regsionsPattern)
	if err != nil {
		return nil, nil, err
	}
	sort.Sort(sort.Reverse(sort.StringSlice(regionDirs)))
	dataMap := map[string][]regionReportData{}
	var dates []time.Time
	for _, d := range regionDirs {
		secsStr := strings.Replace(path.Base(d), ".png", "", 1)
		secs, err := strconv.Atoi(secsStr)
		if err != nil {
			return nil, nil, errors.Errorf("strconv.Atoi: secsStr=%s %v", secsStr, err)
		}
		date := time.Unix(int64(secs), 0)
		dates = append(dates, date)
		for _, r := range regions {
			regionJSONFile := path.Join(d, fmt.Sprintf("%s.json", r))
			regionJSON, err := ioutil.ReadFile(regionJSONFile)
			if err != nil {
				return nil, nil, errors.Errorf("ioutil.ReadFile: regionJSONFile=%s %v", regionJSONFile, err)
			}
			var data []regionRawData
			if err := json.Unmarshal(regionJSON, &data); err != nil {
				// Soft fail
				log.Printf("json.Unmarshal: regionJSONFile=%s %v", regionJSONFile, err)
				continue
			}
			for _, d := range data {
				d.sourceFile = regionJSONFile
			}
			for _, d := range data {
				reportKey := d.toMetadata()
				reportData := d.toData()
				rrd := regionReportData{
					date:       date,
					key:        reportKey,
					data:       reportData,
					sourceFile: regionJSONFile,
				}
				key := fmt.Sprintf("%s:%s", reportKey.state, reportKey.county)
				regionReportDatas, ok := dataMap[key]
				if !ok {
					regionReportDatas = []regionReportData{}
				}
				regionReportDatas = append(regionReportDatas, rrd)
				dataMap[key] = regionReportDatas
			}
		}
	}
	return dataMap, dates, nil
}

func outputRegions(dataDir string, dataMap map[string][]regionReportData, dates []time.Time) error {
	csvDir, err := mkdirAll(dataDir, "csv")
	if err != nil {
		return err
	}

	var sortedKeys []string
	for key := range dataMap {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)

	formatDate := func(d time.Time) string {
		return d.Format("2006-01-02 15:04:05")
	}

	dataForDate := func(datas []regionReportData, date time.Time) regionReportData {
		var closest regionReportData
		var minDur time.Duration
		for _, d := range datas {
			diff := d.date.Sub(date)
			if diff < 0 {
				diff *= -1
			}
			if minDur == 0 || diff < minDur {
				minDur = diff
				closest = d
			}
		}
		return closest
	}

	var sortedDates []time.Time
	sortedDates = append(sortedDates, dates...)
	sort.Slice(sortedDates, func(i, j int) bool {
		return sortedDates[i].Before(sortedDates[j])
	})

	header := []string{"State", "County"}
	for _, d := range sortedDates {
		header = append(header, formatDate(d))
	}

	makeCSVWriter := func(outfile string) (*csv.Writer, error) {
		f := path.Join(csvDir, "percentages.csv")
		log.Printf("writing to %s", f)
		file, err := os.Create(f)
		if err != nil {
			return nil, err
		}
		writer := csv.NewWriter(file)
		return writer, nil
	}

	{
		writer, err := makeCSVWriter("percentages.csv")
		if err != nil {
			return err
		}
		defer writer.Flush()

		if err := writer.Write(header); err != nil {
			return errors.Errorf("writer.Write(header): percentages.csv %v", err)
		}
		for _, k := range sortedKeys {
			datas := dataMap[k]
			line := []string{datas[0].key.state, datas[0].key.county}
			for _, d := range sortedDates {
				data := dataForDate(datas, d)
				p := float64(data.data.customersAffected) / float64(data.data.customersServed)
				line = append(line, fmt.Sprintf("%f", p))
			}
			if err := writer.Write(line); err != nil {
				return errors.Errorf("writer.Write: percentages.csv %v", err)
			}
		}
	}

	{
		writer, err := makeCSVWriter("customers_affected.csv")
		if err != nil {
			return err
		}
		defer writer.Flush()

		if err := writer.Write(header); err != nil {
			return errors.Errorf("writer.Write(header): customers_affected.csv %v", err)
		}
		for _, k := range sortedKeys {
			datas := dataMap[k]
			line := []string{datas[0].key.state, datas[0].key.county}
			for _, d := range sortedDates {
				data := dataForDate(datas, d)
				line = append(line, fmt.Sprintf("%d", data.data.customersAffected))
			}
			if err := writer.Write(line); err != nil {
				return errors.Errorf("writer.Write: customers_affected.csv %v", err)
			}
		}
	}

	return nil
}
