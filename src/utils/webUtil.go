package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func FetchData[T any](url string) (result T, err error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error retrieving Sorting Meta: ", err)
		return
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error unmarshaling response body:", err)
		return
	}
	return
}
