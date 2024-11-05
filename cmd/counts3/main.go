package main

import (
	"context"
	option "counts3/internal/option"
	queue "counts3/internal/queue"
	worker "counts3/internal/worker"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	bucketName, fileName, workerPool := option.Param()
	jobs := make(chan string, 20)

	s3Client := s3Client()

	var wg sync.WaitGroup

	for i := 1; i <= workerPool; i++ {
		wg.Add(1)
		go worker.Worker(s3Client, bucketName, jobs, &wg)
	}

	queue.QueueJob(fileName, jobs)

	close(jobs) // Close the channel to indicate that no more jobs will be added.
	wg.Wait()   // Wait for all workers to finish.
}

func s3Client() *s3.Client {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	return s3.NewFromConfig(cfg)
}
