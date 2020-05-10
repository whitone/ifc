// Copyright 2020 Stefano Cotta Ramusino. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"text/template"
	"time"

	// A small library that turns strings in to slugs
	"github.com/Machiel/slugify"
)

const (
	genFileName      = "municipalities.go"

	municipalURL     = "https://www.istat.it/storage/codici-unita-amministrative/Elenco-comuni-italiani.csv"
	municipalCodeCol = 18
	municipalNameCol = 5

	countriesURL     = "https://www.istat.it/it/files//2011/01/Elenco-codici-e-denominazioni-unita-territoriali-estere.zip"
	countriesFileExt = ".csv"
	countryCodeCol   = 9
	countryNameCol   = 6
)

type municipalities struct {
	ByCode    map[string]string
	ByName    map[string]string
	CreatedOn string
}

func populateMap(m *municipalities, body []byte, codeColumn int, nameColumn int) error {
	utf8Body := make([]rune, len(body))
	for i, b := range body {
		utf8Body[i] = rune(b)
	}
	csvBody := csv.NewReader(strings.NewReader(string(utf8Body)))
	csvBody.Comma = ';'
	_, err := csvBody.Read()
	if err != nil {
		return err
	}
	for {
		csvLine, err := csvBody.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		code := slugify.Slugify(csvLine[codeColumn])
		if code != "" {
			name := csvLine[nameColumn]
			m.ByCode[code] = name
			name = slugify.Slugify(name)
			m.ByName[name] = code
		}
	}
	return nil
}

func body(url string) (body []byte, err error) {
	response, err := http.Get(url)
	if err != nil {
		return
	}
	defer func() {
		if err = response.Body.Close(); err != nil {
			return
		}
	}()
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	return
}

func main() {
	m := &municipalities{}
	m.ByCode = make(map[string]string)
	m.ByName = make(map[string]string)

	b, err := body(municipalURL)
	if err != nil {
		log.Fatal(err)
	}
	err = populateMap(m, b, municipalCodeCol, municipalNameCol)
	if err != nil {
		log.Fatal(err)
	}

	archiveBody, err := body(countriesURL)
	if err != nil {
		log.Fatal(err)
	}
	archive, err := zip.NewReader(bytes.NewReader(archiveBody), int64(len(archiveBody)))
	if err != nil {
		log.Fatal(err)
	}
	b = nil
	for _, archiveFile := range archive.File {
		if path.Ext(archiveFile.Name) == countriesFileExt {
			file, er := archiveFile.Open()
			if er != nil {
				log.Fatal(er)
			}
			defer func() {
				if err := file.Close(); err != nil {
					log.Fatal(err)
				}
			}()
			b, err = ioutil.ReadAll(file)
			if err != nil {
				log.Fatal(err)
			}
			break
		}
	}
	if b == nil {
		log.Fatal(errors.New("no countries csv file found in archive"))
	}
	err = populateMap(m, b, countryCodeCol, countryNameCol)
	if err != nil {
		log.Fatal(err)
	}

	var municipalitiesTemplate = template.Must(template.New("").Parse("" +
		"// File generated on {{ .CreatedOn }} by updateMunicipalities (https://github.com/whitone/ifc)\n\n" +
		"package ifc\n\n" +
		"// municipalitiesByCode contains all Italian municipalities and foreign countries sorted by code.\n" +
		"var municipalitiesByCode = map[string]string{ {{ range $key, $value := .ByCode }}\n" +
		"   \"{{ $key }}\": \"{{ $value }}\",{{ end }}\n" +
		"}\n\n" +
		"// municipalitiesByName contains all Italian municipalities and foreign countries sorted by name.\n" +
		"var municipalitiesByName = map[string]string{ {{ range $key, $value := .ByName }}\n" +
		"   \"{{ $key }}\": \"{{ $value }}\",{{ end }}\n" +
		"}\n"))

	generatedFile, err := os.Create(genFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := generatedFile.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	m.CreatedOn = time.Now().Format(time.RFC1123)
	err = municipalitiesTemplate.Execute(generatedFile, m)
	if err != nil {
		log.Fatal(err)
	}
}
