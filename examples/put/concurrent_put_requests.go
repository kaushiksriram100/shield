// This client program sends several GET requests to server in a concurrent way. User may add -n to specify the concurrency.
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
)

func main() {

	var num_concurrent_requests = flag.Int("n", 500, "number of concurrent requests")
	var url = flag.String("u", "http://localhost:8080/urlinfo/1/google.com:8080/search?v=2", "URL")
	flag.Parse()

	var wg sync.WaitGroup
	for i := 0; i < *num_concurrent_requests; i++ {
		wg.Add(1)
		go func() {
			client := &http.Client{}
			random := rand.Intn(100000)
			randomStr := strconv.FormatInt(int64(random), 10)
			req, err := http.NewRequest("PUT", *url, nil)

			q := req.URL.Query()
			q.Add("v", randomStr)
			req.URL.RawQuery = q.Encode()
			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("Error fetching %s: %v\n", *url, err)
				wg.Done()
				return
			}
			defer resp.Body.Close()

			fmt.Printf("Response status code from %s, %s: %d\n", *url, randomStr, resp.StatusCode)
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println("All requests completed.")
}
