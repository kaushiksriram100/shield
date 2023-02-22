// This client program sends several GET requests to server in a concurrent way. User may add -n to specify the concurrency.
package main

import (
	"flag"
	"fmt"
	"net/http"
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
			resp, err := http.Get(*url)
			if err != nil {
				fmt.Printf("Error fetching %s: %v\n", *url, err)
				wg.Done()
				return
			}
			defer resp.Body.Close()

			fmt.Printf("Response status code from %s: %d\n", *url, resp.StatusCode)
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println("All requests completed.")
}
