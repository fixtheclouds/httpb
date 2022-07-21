package main

import (
	"fmt"
	"net/http"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func doRequest(url string) time.Duration {
	start := time.Now()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while fetching")
	}
	defer resp.Body.Close()
	_, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		fmt.Println(readErr)
	}

	end := time.Now()
	elapsed := end.Sub(start)

	s := fmt.Sprintf("Fetched in %v", elapsed)
	fmt.Println(s)

	return elapsed
}

func printResults(results []time.Duration) {
	s := fmt.Sprintf("Total requests: %d", len(results))
	fmt.Println(s)
}

func main() {
	url := os.Args[1]
	c := make(chan os.Signal)
	var results []time.Duration
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
			<-c
			printResults(results)
			os.Exit(1)
	}()

	for {
		results = append(results, doRequest(url))
	}
}
