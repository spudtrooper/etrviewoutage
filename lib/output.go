package lib

import (
	"github.com/go-errors/errors"
)

func OutputAll(dataDir string) error {
	if err := outputHTML(dataDir); err != nil {
		return errors.Errorf("OutputHTML: %v", err)
	}
	dataMap, dates, err := readRegionDataMap(dataDir)
	if err != nil {
		return errors.Errorf("readRegionDataMap: %v", err)
	}
	if err := outputRegions(dataDir, dataMap, dates); err != nil {
		return errors.Errorf("OutputRegions: %v", err)
	}
	if err := outputCharts(dataDir, dataMap); err != nil {
		return errors.Errorf("OutputCharts: %v", err)
	}
	return nil
}
