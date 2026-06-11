package crawler

import (
	"fmt"
	"net/http"
	"time"
)

var client = &http.Client{
	Timeout: 10 * time.Second,
}

func Crawl(url string) {
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("Error fetching %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf(
		"Fetched %s (status: %d)\n",
		url,
		resp.StatusCode,
	)
}
