package api

import (
	"log"
	"fmt"
	"net/url"
	"net/http"
	"encoding/json"
	"net/http/httputil"
)

const QUERY_STRING = "select temp_high, temp_low, temp_med, temp_coll from puffer where time > now() - 1h order desc limit 1"

type column string
type point int

type QueryResult struct {
	Name    string
	Columns []column
	Points  [][]point
}

// FetchPufferData is for getting the latest puffer data
func FetchPufferData(options *Options) (*Info, error) {

	baseUrl := options.Url
	if baseUrl == "" {
		return nil, fmt.Errorf("No influxdb URL provided")
	}
	url :=
		fmt.Sprintf(options.Url+"?u=%s&p=%s&q=%s",
			url.QueryEscape(options.User), url.QueryEscape(options.Password), url.QueryEscape(QUERY_STRING))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	data := make([]QueryResult, 0)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	dump, _ := httputil.DumpResponse(resp, true)

	log.Print(string(dump[:]))
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	result := data[0]

	fmt.Printf("%+v\n", result)
	var collector, high, med, low float32
	for i := 0; i < len(result.Columns); i++ {
		key := result.Columns[i]
		value := result.Points[0][i]
		if key == "temp_coll" {
			collector = float32(value / 10)
		} else if key == "temp_high" {
			high = float32(value / 10)
		} else if key == "temp_med" {
			med = float32(value / 10)
		} else if key == "temp_low" {
			low = float32(value / 10)
		}
	}

	log.Printf("High: %f -- Med: %f -- Low: %f -- Collector: %f", high, med, low, collector)

	return &Info{
		CollectorTemp: collector,
		HighTemp:      high,
		LowTemp:       low,
		MidTemp:       med,
	}, nil
}
