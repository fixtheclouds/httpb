package main

import (
	"fmt"
	"errors"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Result struct {
	Status string
	Time   int
}

func doRequest(url string) (Result, error) {
	start := time.Now()

	resp, err := http.Get(url)
	if err != nil {
		return Result{}, errors.New("Failed to fetch URL")
	}
	_, readErr := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if readErr != nil {
		return Result{}, errors.New("Failed to parse body")
	}

	end := time.Now()
	elapsed := end.Sub(start)

	s := fmt.Sprintf("Fetched in %v, status: %s", elapsed.Round(time.Millisecond), resp.Status)
	fmt.Println(s)

	return Result{resp.Status, toMs(elapsed)}, nil
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

func printResults(results []Result, failuresCount int) {
	if len(results) == 0 {
		fmt.Println("No successful results")
		return
	}

	min, avg, max := findMinAvgMax(results)
	s := fmt.Sprintf("Min: %v ms\nAvg: %v ms\nMax: %v ms", min, avg, max)
	fmt.Println(s)
	s = fmt.Sprintf("Successful requests: %d\nFailed requests: %d", len(results), failuresCount)
	fmt.Println(s)
}

func toMs(duration time.Duration) int {
	return int(duration / time.Millisecond)
}

func main() {
	url := os.Args[1]
	c := make(chan os.Signal)
	failuresCount := 0
	var results []Result
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		printResults(results, failuresCount)
		os.Exit(1)
	}()

	for {
		result, error := doRequest(url)
		if error != nil {
			fmt.Println(error)
			failuresCount += 1
		} else {
			results = append(results, result)
		}
	}
}
