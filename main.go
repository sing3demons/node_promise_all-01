package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type RequestData struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func createData(data RequestData) ([]byte, error) {
	requestBody, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	// Make an HTTP POST request
	resp, err := http.Post("http://localhost:3000/example", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func createSyncData(data []RequestData) []RequestData {
	var response []RequestData
	for _, body := range data {
		result, err := createData(body)
		if err != nil {
			fmt.Println("Error:", err)
			return nil
		}
		data := RequestData{}
		json.Unmarshal(result, &data)
		response = append(response, data)
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
		go func(data RequestData) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			responseData, err := httpGet(data)
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

func main() {
	data := []RequestData{}
	for i := 0; i < 900; i++ {
		data = append(data, RequestData{Name: fmt.Sprintf("Name_%d", i)})
	}
	start := time.Now()
	// result := createSyncData(data) // Time: 10.354932109s
	result := createAsyncData(data) // Time: 5.04395879s
	for _, data := range result {
		fmt.Println(data)
	}
	end := time.Now()
	fmt.Printf("Time: %s\n", end.Sub(start))
}

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

func httpGetWithRetry(data RequestData, maxRetries int) (*RequestData, error) {
	timeout := 30 * time.Second
	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		httpReq, err := http.NewRequest(http.MethodPost, "http://localhost:3000/example", bytes.NewBuffer(b))
		if err != nil {
			return nil, err
		}
		httpReq = httpReq.WithContext(ctx)
		httpReq.Header.Set("Content-Type", "application/json")
		httpClient := &http.Client{
			Timeout: timeout,
		}
		resp, err := httpClient.Do(httpReq)
		if err != nil {
			// Log the error
			fmt.Println("Attempt", attempt, "- Error:", err)
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			// Log the error
			fmt.Println("Attempt", attempt, "- Error:", err)
			continue
		}

		var result RequestData
		err = json.Unmarshal(body, &result)
		if err != nil {
			return nil, err
		}

		return &result, nil
	}

	return nil, fmt.Errorf("max retries exceeded")
}
