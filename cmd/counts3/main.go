package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/alvinsiew/counts3/internal/worker"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func worker(s3c *s3.Client, bucket string, jobs chan string, wg *sync.WaitGroup) {

	for job := range jobs {
		countFilesInS3Folder(s3c, bucket, job)
	}
	defer wg.Done()
}

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
		go worker(s3Client, *bucketParam, jobs, &wg)
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

// countFilesInS3Folder counts the number of objects in a specified S3 folder
func countFilesInS3Folder(client *s3.Client, bucket string, prefix string) {
	var count int

	for {
		input := &s3.ListObjectsV2Input{
			Bucket: aws.String(bucket),
			Prefix: aws.String(prefix),
		}

		result, err := client.ListObjectsV2(context.TODO(), input)
		if err != nil {
			log.Fatal(err)
		}

		for _, object := range result.Contents {
			// Exclude the folder itself (usually represented as an object with size 0)
			if *object.Size > 0 {
				count++
			}
		}

		if !*result.IsTruncated {
			break // No more objects to retrieve
		}
	}

	fmt.Printf("Total files in folder %s: %d\n", prefix, count)
}
