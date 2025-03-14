package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"CheckSorting/types"
	"CheckSorting/utils"
)

type intensiveCalculationMeta struct {
	SortType         string
	AverageTimeTaken types.ISO8601Duration
	MinTimeTaken     types.ISO8601Duration
	MaxTimeTaken     types.ISO8601Duration
	SampleSize       uint64
}

func createLongCalculations(nRows int, nColumns int) [][]int {
	rows := make([][]int, 0)
	for range nRows {
		row := make([]int, 0)
		for range nColumns {
			row = append(row, rand.Int())
		}
		rows = append(rows, row)
	}
	return rows
}

func sendSortingRequest(url string, nRows int, nColumns int) {
	requestItem := createLongCalculations(nRows, nColumns)
	jsonData, err := json.Marshal(requestItem)
	if err != nil {
		fmt.Println("Error marshalling JSON: ", err)
		return
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("HTTP request failed: ", err)
		return
	}
	defer resp.Body.Close()
}

func retrieveSortingMeta(url string, nRequestedSorts int) {
	for {
		parsedResp, err := utils.FetchData[[]intensiveCalculationMeta](url)
		if err != nil {
			return
		}
		time.Sleep(1 * time.Second)
		if len(parsedResp) == 2 {
			ok := true
			for _, item := range parsedResp {
				if item.SampleSize < uint64(nRequestedSorts) {
					ok = false
				}
			}
			if ok {
				for _, item := range parsedResp {
					fmt.Println("SortType", item.SortType)
					fmt.Println("MinTimeTaken", item.MinTimeTaken)
					fmt.Println("AverageTimeTaken", item.AverageTimeTaken)
					fmt.Println("MaxTimeTaken", item.MaxTimeTaken)
					fmt.Println()
				}
				return
			}
		}
	}
}

func deletePreviousData(url string) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()
}

func main() {
	url := "http://127.0.0.1:4041/sort"
	urlDeleteAll := url + "/delete-all"
	urlRetrieveMeta := url + "/meta"
	urlQuick := url + "/calculative/calculate-once"
	urlSlow := url + "/calculative/intensive"
	nRequests := 10
	nRows := 512
	nColumns := 512

	deletePreviousData(urlDeleteAll)
	for range nRequests {
		go sendSortingRequest(urlQuick, nRows, nColumns)
		go sendSortingRequest(urlSlow, nRows, nColumns)
	}

	retrieveSortingMeta(urlRetrieveMeta, nRequests)
}
