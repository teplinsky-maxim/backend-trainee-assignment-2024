package main

import (
	"avito-backend-2024-trainee/benchmarks/prepare"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"net/http"
	"sync"
	"time"
)

// Pure DB
//Average Response Time: 326.665456ms
//Minimum Response Time: 53.4741ms
//Maximum Response Time: 557.6233ms

// With in-memory cache
//Average Response Time (per second): 614.0676ms
//Minimum Response Time (per second): 322.0484ms
//Maximum Response Time (per second): 1.2736s

func main() {
	//prepare.FillDatabase()
	Bench()
}

func Bench() {
	url := "http://localhost:3000/api/v1/user_banner"
	token := "user"

	// Variables to hold response time statistics
	var totalResponseTime time.Duration
	var minResponseTime, maxResponseTime time.Duration
	minResponseTime = time.Second * 10 // Set initial value to a high number

	requestsPerSecond := 200

	// Create a ticker to send requests every second
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	wg := sync.WaitGroup{}
	// Loop to send requests
	for range ticker.C {
		for i := 0; i < requestsPerSecond; i++ {
			i := i
			wg.Add(1)
			go func() {
				sendSingleRequest(url, token, i, &totalResponseTime, &minResponseTime, &maxResponseTime)
				wg.Done()
			}()
		}

		wg.Wait()

		printStatistics(totalResponseTime, minResponseTime, maxResponseTime, requestsPerSecond)

		totalResponseTime = 0
		minResponseTime = time.Second * 10
		maxResponseTime = 0
	}
}

func sendSingleRequest(url, token string, i int, totalResponseTime, minResponseTime, maxResponseTime *time.Duration) {
	var tagId, featureId uint
	if i%4 == 0 {
		featureId = prepare.CommonFeature
	} else {
		// TODO: const для значений
		featureId = gofakeit.UintRange(1, 1000)
	}

	if i%3 == 0 {
		tagId = prepare.CommonTag1
	} else if i%4 == 0 {
		tagId = prepare.CommonTag2
	} else {
		tagId = gofakeit.UintRange(1, 1000)
	}
	params := fmt.Sprintf("tag_id=%v&feature_id=%v&use_last_revision=false", tagId, featureId)
	startTime := time.Now()

	// Create request with parameters
	req, err := http.NewRequest("GET", url+"?"+params, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Add Authorization header
	req.Header.Set("Authorization", "Bearer "+token)

	// Send request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	responseTime := time.Since(startTime)

	// Update statistics
	if responseTime < *minResponseTime {
		*minResponseTime = responseTime
	}
	if responseTime > *maxResponseTime {
		*maxResponseTime = responseTime
	}
	*totalResponseTime += responseTime
}

func printStatistics(totalResponseTime, minResponseTime, maxResponseTime time.Duration, numRequests int) {
	avgResponseTime := totalResponseTime / time.Duration(numRequests)
	fmt.Println("Average Response Time (per second):", avgResponseTime)
	fmt.Println("Minimum Response Time (per second):", minResponseTime)
	fmt.Println("Maximum Response Time (per second):", maxResponseTime)
	fmt.Println("")
}
