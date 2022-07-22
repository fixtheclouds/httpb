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

func doRequest(url string) (time.Duration) {
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

	s := fmt.Sprintf("Fetched in %v", elapsed.Round(time.Millisecond))
	fmt.Println(s)

	return elapsed
}

func findMinAvgMax(results []int) (min int, avg float64, max int) {
	sum := 0
	min = results[0]
	max = results[0]
	for _, value := range results {
		sum += value
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}
	avg = float64(sum) / float64(len(results))

	return min, avg, max
}

func printResults(results []int) {
	min, avg, max := findMinAvgMax(results)
	s := fmt.Sprintf("Min: %v ms\nAvg: %v ms\nMax: %v ms\nTotal requests: %d", min, avg, max, len(results))
	fmt.Println(s)
}

func toMs(duration time.Duration) (int) {
	return int(duration / time.Millisecond)
}

func main() {
	url := os.Args[1]
	c := make(chan os.Signal)
	var results []int
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
			<-c
			printResults(results)
			os.Exit(1)
	}()

	for {
		results = append(results, toMs(doRequest(url)))
	}
}
