package server

import (
	"crypto/tls"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	chart "github.com/wcharczuk/go-chart"
)

func addMdCmd(url string) string {
	if strings.ContainsRune(url, '?') {
		return url + "md()"
	}
	return url + "?md()"
}

func handlePlotForm(w http.ResponseWriter, r *http.Request) {
	t, _ := template.New("plotIndex").Parse(plotIndex())
	t.Execute(w, nil)
}

func retrieveData(dataurl string) (string, error) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	//client := &http.Client{}
	httpreq, err := http.NewRequest(http.MethodGet, addMdCmd(dataurl), nil)
	if err != nil {
		return "", err
	}

	httpresp, err := client.Do(httpreq)
	if err != nil {
		return "", err
	}

	if httpresp.StatusCode != http.StatusOK {
		log.Print(httpresp.StatusCode)
		return "", errors.New(httpresp.Status)
	}

	respbody, err := ioutil.ReadAll(httpresp.Body)
	if err != nil {
		fmt.Println("(2)")
		return "", err
	}

	return string(respbody), nil
}

func findTimeColumn(header []string) (int, error) {
	for x := range header {
		if strings.Contains(header[x], "{dc:date}") ||
			strings.Contains(header[x], "{yamz:h1317}") {
			return x, nil
		}
	}
	return -1, errors.New("Time attribute not found")
}

func dataTime(data []string) []time.Time {
	var tdata []time.Time
	for x := range data {
		if data[x] == "" {
			continue
		}
		t, err := time.Parse("2006-01-02 15:04:05", data[x])
		if err != nil {
			log.Print(err)
		}
		tdata = append(tdata, t)
	}
	return tdata
}

func dataFloat(data []string) []float64 {
	var fdata []float64
	for x := range data {
		if data[x] == "" {
			data[x] = "0"
		}
		f, _ := strconv.ParseFloat(data[x], 64)
		fdata = append(fdata, f)
	}
	return fdata
}

func makeSeries(dataurl string) []chart.Series {

	rawdata, _ := retrieveData(dataurl)

	var series []chart.Series

	datarows := strings.Split(rawdata, "\n")
	header := strings.Split(datarows[0], ",")

	timeIndex, err := findTimeColumn(header)
	if err != nil {
		log.Print(err)
		// TODO Handle error.
	}

	var d [][]string
	d = make([][]string, len(header))
	for x := range d {
		d[x] = make([]string, 1000000)
	}

	for r := range datarows {
		if r == 0 {
			continue
		}
		row := strings.Split(datarows[r], ",")
		for c := range row {
			d[c][r-1] = row[c]
		}
	}

	for c := range header {
		if c == timeIndex {
			continue
		}
		h := header[c]
		x := strings.IndexRune(h, '{')
		if x >= 0 {
			h = h[0:x]
		}

		series = append(series,
			chart.TimeSeries{
				Name:    h,
				XValues: dataTime(d[timeIndex]),
				YValues: dataFloat(d[c]),
			})
	}

	return series
}

func handlePlotRun(w http.ResponseWriter, r *http.Request) {

	graph := chart.Chart{
		XAxis: chart.XAxis{
			Style: chart.Style{Show: true},
		},
		YAxis: chart.YAxis{
			Style: chart.Style{Show: true},
		},
		Background: chart.Style{
			Padding: chart.Box{
				Top:  20,
				Left: 20,
			},
		},
		Series: makeSeries(r.URL.Query()["dataurl"][0]),
	}

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	w.Header().Set("Content-Type", "image/png")
	_ = graph.Render(chart.PNG, w)
}

func handlePlot(w http.ResponseWriter, r *http.Request) {
	if r.URL.RawQuery == "" {
		handlePlotForm(w, r)
	} else {
		handlePlotRun(w, r)
	}
}
