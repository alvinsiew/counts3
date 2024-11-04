package main

import (
	"bufio"
	"context"
	worker "counts3/internal/worker"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	var filePath string

	fileParam := flag.String("f", "", "File path and name")
	workerParam := flag.Int("w", 1, "Number of Workers")
	bucketParam := flag.String("s3", "", "S3 bucket name")
	flag.Parse()

	if *fileParam == "" {
		fmt.Println("File must be specific with -f")
		os.Exit(1)
	} else {
		filePath = *fileParam
	}

	if *bucketParam == "" {
		fmt.Println("S3 Bucket must be specific with -s3")
		os.Exit(1)
	}

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	var s3Client = s3.NewFromConfig(cfg)

	jobs := make(chan string, 100)
	var wg sync.WaitGroup

	for i := 1; i <= *workerParam; i++ {
		wg.Add(1)
		go worker.Worker(s3Client, *bucketParam, jobs, &wg)
	}

	file, err := os.Open(filePath)
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
