package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	version = "2.2.0"
	fmtURL  = "https://reqrypt.org/download/WinDivert-%s-A.zip"
)

func main() {
	url := fmt.Sprintf(fmtURL, version)
	response, err := http.Get(url)
	if err != nil || response.StatusCode != http.StatusOK {
		log.Fatal("fail to download DLL")
	}
	defer response.Body.Close()
	zipFile, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	reader := bytes.NewReader(zipFile)
	zipReader, err := zip.NewReader(reader, int64(reader.Len()))
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range zipReader.File {
		if strings.HasSuffix(file.Name, "x64/WinDivert.dll") {
			saveFile(file)
		}
		if strings.HasSuffix(file.Name, "x64/WinDivert.lib") {
			saveFile(file)
		}
		if strings.HasSuffix(file.Name, "x64/WinDivert64.sys") {
			saveFile(file)
		}
		if strings.HasSuffix(file.Name, "x86/WinDivert.dll") {
			saveFile(file)
		}
		if strings.HasSuffix(file.Name, "x86/WinDivert.lib") {
			saveFile(file)
		}
		if strings.HasSuffix(file.Name, "x86/WinDivert32.sys") {
			saveFile(file)
		}
		if strings.HasSuffix(file.Name, "x86/WinDivert64.sys") {
			saveFile(file)
		}
	}
}

func saveFile(file *zip.File) {
	var filename string
	vpath := fmt.Sprintf("WinDivert-%s-A", version)
	f, err := file.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	switch file.Name {
	case filepath.ToSlash(filepath.Join(vpath, "x64", filepath.Base(file.Name))):
		filename = filepath.Join("x64", filepath.Base(file.Name))
		err = os.MkdirAll("x64", 0755)
		if err != nil {
			log.Fatal(err)
		}
	case filepath.ToSlash(filepath.Join(vpath, "x86", filepath.Base(file.Name))):
		filename = filepath.Join("x86", filepath.Base(file.Name))
		err = os.MkdirAll("x86", 0755)
		if err != nil {
			log.Fatal(err)
		}
	default:
		filename = filepath.Base(file.Name)
	}
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
