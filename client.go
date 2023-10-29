package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

func httpGet(data RequestData) (*RequestData, error) {
	// responseCh chan RequestData
	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	httpReq, err := http.NewRequest(http.MethodPost, "http://localhost:3000/example", bytes.NewBuffer(b))
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	var result RequestData
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func createData(data []RequestData) []RequestData {
	var response []RequestData
	for _, body := range data {
		result, err := httpGet(body)
		if err != nil {
			fmt.Println("Error:", err)
			return nil
		}
		response = append(response, *result)
	}
	return response
}

func createAsyncData(data []RequestData) []RequestData {
	var response []RequestData
	var wg sync.WaitGroup
	responseCh := make(chan RequestData, len(data))
	poolSize := 10
	semaphore := make(chan struct{}, poolSize)
	for i := 0; i < len(data); i++ {
		wg.Add(1)
		go func(body RequestData) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			responseData, err := httpGet(body)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			responseCh <- *responseData

		}(data[i])
	}

	go func() {
		wg.Wait()
		close(responseCh)
	}()
	for data := range responseCh {
		response = append(response, data)
	}

	return response
}
