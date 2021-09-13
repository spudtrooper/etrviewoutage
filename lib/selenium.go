package lib

import (
	"bytes"
	"fmt"
	"image/png"
	"log"
	"os"

	"github.com/go-errors/errors"
	"github.com/spudtrooper/goutil/or"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

func findChromeDriver() string {
	fileExists := func(f string) bool {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			return false
		}
		return true
	}
	paths := []string{
		"/opt/homebrew/bin/chromedriver",
		"/usr/local/bin/chromedriver",
	}
	for _, f := range paths {
		if fileExists(f) {
			return f
		}
	}
	return ""
}

type makeWebDriverOptions struct {
	Verbose          bool
	Headless         bool
	ChromeDriverPath string
	SeleniumPath     string
	Port             int
}

func makeWebDriver(opts makeWebDriverOptions) (selenium.WebDriver, func(), error) {
	seleniumPath := or.String(opts.SeleniumPath, "third_party/selenium/vendor/selenium-server.jar")
	port := or.Int(opts.Port, 8082)
	chromeDriverPath := or.String(opts.ChromeDriverPath, findChromeDriver())
	if chromeDriverPath == "" {
		return nil, nil, errors.Errorf("Couldn't find chromedriver")
	}
	selOpts := []selenium.ServiceOption{
		selenium.ChromeDriver(chromeDriverPath),
	}
	if opts.Verbose {
		selOpts = append(selOpts, selenium.Output(os.Stderr))
		selenium.SetDebug(true)
	}
	service, err := selenium.NewSeleniumService(seleniumPath, port, selOpts...)
	if err != nil {
		return nil, nil, err
	}

	args := []string{
		"--no-sandbox",
		"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/604.4.7 (KHTML, like Gecko) Version/11.0.2 Safari/604.4.7",
	}
	if opts.Headless {
		args = append(args, "--headless")
	}
	caps := selenium.Capabilities{"browserName": "chrome"}
	chromeCaps := chrome.Capabilities{
		Path: "",
		Args: args,
	}
	caps.AddChrome(chromeCaps)
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		return nil, nil, err
	}

	return wd, func() {
		wd.Quit()
		service.Stop()
	}, nil
}

func findElement(wd selenium.WebDriver, tag, text string) (selenium.WebElement, error) {
	els, err := wd.FindElements(selenium.ByTagName, tag)
	if err != nil {
		return nil, err
	}
	for _, el := range els {
		txt, err := el.Text()
		if err != nil {
			return nil, err
		}
		if txt == text {
			return el, nil
		}
	}
	return nil, nil
}

func waitForElement(wd selenium.WebDriver, tagName, text string) (selenium.WebElement, error) {
	var res selenium.WebElement
	var cnt int
	wd.Wait(func(wd selenium.WebDriver) (bool, error) {
		log.Printf("waiting for div %s [%d] ...", text, cnt+1)
		cnt++
		btn, err := findElement(wd, tagName, text)
		if err != nil {
			return false, err
		}
		if btn != nil {
			res = btn
			return true, nil
		}
		return false, nil
	})
	if res == nil {
		return nil, fmt.Errorf("couldn't find div with text: %s", text)
	}
	return res, nil
}

// From: https://github.com/lucasmdomingues/go-selenium-screenshot/blob/master/print.go
func takeScreenshot(wd selenium.WebDriver, outFile string) error {
	ss, err := wd.Screenshot()
	if err != nil {
		return err
	}
	r := bytes.NewReader(ss)

	im, err := png.Decode(r)
	if err != nil {
		return err
	}

	log.Printf("writing to %s...", outFile)
	f, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	png.Encode(f, im)

	return nil
}
