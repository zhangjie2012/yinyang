package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

var (
	SourceUrl = "https://www.hko.gov.hk/tc/gts/time/calendar/text/files/T%dc.txt"
	YearStart = 1901
	YearEnd   = 2101
)

func Decodebig5(s []byte) ([]byte, error) {
	I := bytes.NewReader(s)
	O := transform.NewReader(I, traditionalchinese.Big5.NewDecoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func httpGet(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	rawBs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	bs, err := Decodebig5(rawBs)
	if err != nil {
		log.Fatal(err)
	}

	return string(bs)
}

func Crawler(year int) (data []string) {
	data = make([]string, 0)
	url := fmt.Sprintf(SourceUrl, year)
	content := httpGet(url)
	rows := strings.Split(content, "\n")
	for idx, row := range rows {
		if idx < 3 {
			continue
		}
		fields := strings.Fields(row)
		if len(fields) < 3 {
			continue
		}
		data = append(data, strings.Join(fields, ","))
	}
	return data
}

func WriteFile(fn string, data []string) {
	f, err := os.Create(fn)
	if err != nil {
		log.Fatal(f)
	}

	for _, d := range data {
		_, err := f.WriteString(fmt.Sprintf("%s\n", d))
		if err != nil {
			log.Fatal(err)
		}
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	for y := YearStart; y < YearEnd; y++ {
		d := Crawler(y)
		fn := fmt.Sprintf("../../rawdata/%d.txt", y)
		WriteFile(fn, d)
		log.Printf("-------------- success, %d\n", y)
	}
}
