package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func outputHTML(dataDir string) error {
	screenshotsDir := path.Join(dataDir, "screenshots")
	screenshotsPattern := path.Join(screenshotsDir, "*.png")
	imagePaths, err := filepath.Glob(screenshotsPattern)
	if err != nil {
		return err
	}
	sort.Strings(imagePaths)

	// Only keep images on or after 1630465902.
	var filteredImagePaths []string
	for _, imgPath := range imagePaths {
		idStr := strings.Replace(path.Base(imgPath), ".png", "", 1)
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return err
		}
		if id >= 1630465902 {
			filteredImagePaths = append(filteredImagePaths, imgPath)
		}
	}

	var images []string
	for _, imgPath := range filteredImagePaths {
		img := fmt.Sprintf("../screenshots/%s", path.Base(imgPath))
		images = append(images, img)
	}
	htmlDir, err := mkdirAll(dataDir, "html")
	if err != nil {
		return err
	}

	animateHtmlIn := path.Join("static", "animate.html")
	animateHtmlOut := path.Join(htmlDir, "animate.html")
	if err := copyFile(animateHtmlIn, animateHtmlOut); err != nil {
		return err
	}

	animateImagesJSON, err := json.Marshal(images)
	if err != nil {
		return err
	}
	animateImagesOut := path.Join(htmlDir, "animate_images.json")
	if err := writeFile(animateImagesOut, animateImagesJSON); err != nil {
		return err
	}

	return nil
}

// https://opensource.com/article/18/6/copying-files-go
func copy(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	if _, err := io.Copy(destination, source); err != nil {
		return err
	}

	return nil
}
