package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

// Job represents your custom job struct.
type Job struct {
	ID   int    `json:"id"`
	Data string `json:"data"`
}

var jobChan chan Job

func handlePostRequest(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var newJob Job
	if err := decoder.Decode(&newJob); err != nil {
		http.Error(rw, "Invalid JSON data", http.StatusBadRequest)
		return
	}
	// Add the new job to your queue (e.g., push it to a channel)
	jobChan <- newJob
	log.Printf("Received new job: %+v", newJob)
}

func worker(ctx context.Context, jobChan <-chan Job) {
	for {
		select {
		case <-ctx.Done():
			return // Worker exits when context is canceled
		case job, ok := <-jobChan:
			if !ok {
				// Channel is closed, no more jobs. Worker exits.
				return
			}
			// Process the job
			processJob(job)
		}
	}
}

func processJob(job Job) {
	// Your job processing logic here
	log.Printf("Processing job: %+v", job)
}

func main() {
	const MaxQueueSize = 5
	jobChan = make(chan Job, MaxQueueSize)
	// Start worker goroutines
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(ctx, jobChan)
		}()
	}

	// Set up HTTP server to handle POST requests
	http.HandleFunc("/add-job", handlePostRequest)
	log.Println("Listening on :8080...")
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// Add example jobs to the queue (optional)
	for i := 0; i < 10; i++ {
		job := Job{ID: i, Data: "Some data"}
		jobChan <- job
	}

	// Wait for workers to finish and cancel context
	wg.Wait()
	cancel()
}
