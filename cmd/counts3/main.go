package main

import (
	"bufio"
	"context"
	option "counts3/internal/option"
	worker "counts3/internal/worker"
	"log"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	bucketName, fileName, workerPool := option.Param()

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	var s3Client = s3.NewFromConfig(cfg)

	// Make channel for queuing tasks
	jobs := make(chan string, 100)

	var wg sync.WaitGroup

	for i := 1; i <= workerPool; i++ {
		wg.Add(1)
		go worker.Worker(s3Client, bucketName, jobs, &wg)
	}

	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		job := scanner.Text()
		jobs <- job
	}

	close(jobs) // Close the channel to indicate that no more jobs will be added.
	wg.Wait()   // Wait for all workers to finish.
}
