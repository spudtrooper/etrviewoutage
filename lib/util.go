package lib

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"text/template"
	"time"
)

var (
	verboseIO = flag.Bool("verboseio", false, "Print verbose IO messages")
)

func mkdirAll(paths ...string) (string, error) {
	outDir := path.Join(paths...)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return "", err
	}
	return outDir, nil
}

func writeFile(f string, b []byte) error {
	if *verboseIO {
		log.Printf("writing to %s", f)
	}
	if err := ioutil.WriteFile(f, b, 0755); err != nil {
		return err
	}
	return nil
}

func renderTemplate(t string, name string, data interface{}) (string, error) {
	tmpl, err := template.New(name).Parse(strings.TrimSpace(t))
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func lastUpdated() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func copyFile(src, dst string) error {
	log.Printf("copy %s -> %s", src, dst)
	fin, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fin.Close()

	fout, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer fout.Close()

	_, err = io.Copy(fout, fin)

	return err
}

func formatPercentage(f float64) string {
	return fmt.Sprintf("%.1f%%", 100*f)
}
