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

type Result struct {
	Status string
	Time int
}

func doRequest(url string) (Result) {
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

	s := fmt.Sprintf("Fetched in %v, status: %s", elapsed.Round(time.Millisecond), resp.Status)
	fmt.Println(s)

	return Result{resp.Status, toMs(elapsed)}
}

func findMinAvgMax(results []Result) (int, float64, int) {
	sum := 0
	min := results[0].Time
	max := results[0].Time
	for _, value := range results {
		sum += value.Time
		if value.Time < min {
			min = value.Time
		}
		if value.Time > max {
			max = value.Time
		}
	}
	avg := float64(sum) / float64(len(results))

	return min, avg, max
}

func printResults(results []Result) {
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
	var results []Result
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
