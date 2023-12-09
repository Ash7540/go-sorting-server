package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"
)

type Payload struct {
	ToSort [][]int `json:"to_sort"`
}

type ProcessResponse struct {
	SortedArrays [][]int `json:"sorted_arrays"`
	TimeNs       int64   `json:"time_ns"`
}

func processingSingleSort(payloader Payload) ProcessResponse {
	startTime := time.Now()
	sortedArrays := make([][]int, len(payloader.ToSort))

	for single_sort := range payloader.ToSort {
		sort.Ints(payloader.ToSort[single_sort])
		sortedArrays[single_sort] = payloader.ToSort[single_sort]
	}

	timeTaken := time.Since(startTime).Nanoseconds()
	return ProcessResponse{
		SortedArrays: sortedArrays,
		TimeNs:       timeTaken,
	}
}

func processingConcurrentSort(payloader Payload) ProcessResponse {
	startTime := time.Now()
	var wg sync.WaitGroup
	result := make([][]int, len(payloader.ToSort))
	ch := make(chan int, len(payloader.ToSort))

	for concurrent_sort, sub_Array := range payloader.ToSort {
		wg.Add(1)
		go func(index int, arr []int) {
			defer wg.Done()
			sort.Ints(arr)
			result[index] = arr
			ch <- index
		}(concurrent_sort, sub_Array)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for range ch {
		// receiving from the channel to confirm that all sorting procedures have been finished
	}

	timeTaken := time.Since(startTime).Nanoseconds()
	return ProcessResponse{
		SortedArrays: result,
		TimeNs:       timeTaken,
	}
}

func processSingleSortHandler(w http.ResponseWriter, r *http.Request) {
	var payloader Payload
	errors := json.NewDecoder(r.Body).Decode(&payloader)
	if errors != nil {
		http.Error(w, errors.Error(), http.StatusBadRequest)
		return
	}

	response_taken_to_sort := processingSingleSort(payloader)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response_taken_to_sort)
}

func processConcurrentSortHandler(w http.ResponseWriter, r *http.Request) {
	var payloader Payload
	errors := json.NewDecoder(r.Body).Decode(&payloader)
	if errors != nil {
		http.Error(w, errors.Error(), http.StatusBadRequest)
		return
	}

	response := processingConcurrentSort(payloader)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/process-single", processSingleSortHandler)
	http.HandleFunc("/process-concurrent", processConcurrentSortHandler)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
