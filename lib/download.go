package lib

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/tebeka/selenium"
)

func Download(outDir string, verbose, headless bool, pause time.Duration) error {
	id := fmt.Sprintf("%d", time.Now().Unix())
	if err := screenshot(outDir, id, verbose, headless, pause); err != nil {
		return err
	}
	if err := downloadJSON(outDir, id); err != nil {
		return err
	}

	return nil
}

func downloadJSON(outDir string, id string) error {
	outDir, err := mkdirAll(outDir, "regions", id)
	if err != nil {
		return err
	}
	for _, s := range regions {
		url := fmt.Sprintf("https://entergy.datacapable.com/datacapable/v1/entergy/Entergy%s/county", s)
		outFile := path.Join(outDir, s+".json")
		log.Printf("downloading %s -> %s", url, outFile)
		if err := downloadFile(outFile, url); err != nil {
			return err
		}
	}
	return nil
}

// https://golangcode.com/download-a-file-from-a-url/
func downloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func screenshot(outDir string, id string, verbose, headless bool, pause time.Duration) error {
	const (
		url = "https://www.etrviewoutage.com/map?state=nola&_ga=2.56165483.1161628684.1630324617-313625520.1630324617"
	)

	outDir, err := mkdirAll(outDir, "screenshots")
	if err != nil {
		return err
	}

	wd, cancel, err := makeWebDriver(makeWebDriverOptions{
		Verbose:  verbose,
		Headless: headless,
	})
	if err != nil {
		return err
	}
	defer cancel()

	log.Printf("screenshotting %s...", url)
	if err := wd.Get(url); err != nil {
		return err
	}

	var waitCnt int
	wd.Wait(func(wd selenium.WebDriver) (bool, error) {
		log.Printf("waiting for logo (%d)...", waitCnt)
		waitCnt++
		imgs, err := wd.FindElements(selenium.ByID, "entergy-map-logo")
		if err != nil {
			return false, err
		}
		if len(imgs) >= 1 {
			return true, nil
		}
		return false, nil
	})

	log.Printf("waiting %v...", pause)
	time.Sleep(pause)

	btn, err := waitForElement(wd, "span", "CLOSE MENU")
	if err != nil {
		return err
	}
	if err := btn.Click(); err != nil {
		return err
	}
	log.Printf("waiting %v...", pause)
	time.Sleep(pause)

	outFile := path.Join(outDir, fmt.Sprintf("%s.png", id))

	if err := takeScreenshot(wd, outFile); err != nil {
		return err
	}

	return nil
}
