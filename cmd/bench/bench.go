package main

import (
	"avito-backend-2024-trainee/benchmarks/prepare"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"net/http"
	"sync"
	"time"
)

//Average Response Time: 326.665456ms
//Minimum Response Time: 53.4741ms
//Maximum Response Time: 557.6233ms

func main() {
	//prepare.FillDatabase()
	Bench()
}

func Bench() {
	url := "http://localhost:3000/api/v1/user_banner"
	numRequests := 200
	token := "user"

	// Variables to hold response time statistics
	var totalResponseTime time.Duration
	var minResponseTime, maxResponseTime time.Duration
	minResponseTime = time.Second * 10 // Set initial value to a high number

	var wg sync.WaitGroup
	wg.Add(numRequests)

	for i := 0; i < numRequests; i++ {
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

		go func() {
			defer wg.Done()

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
			if responseTime < minResponseTime {
				minResponseTime = responseTime
			}
			if responseTime > maxResponseTime {
				maxResponseTime = responseTime
			}
			totalResponseTime += responseTime
		}()
	}

	wg.Wait()

	// Calculate average response time
	avgResponseTime := totalResponseTime / time.Duration(numRequests)

	// Print statistics
	fmt.Println("Average Response Time:", avgResponseTime)
	fmt.Println("Minimum Response Time:", minResponseTime)
	fmt.Println("Maximum Response Time:", maxResponseTime)
}
