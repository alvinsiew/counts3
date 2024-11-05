package worker

import (
	count "counts3/internal/count"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Create worker pool to execute job
func Worker(s3c *s3.Client, bucket string, jobs chan string, wg *sync.WaitGroup) {

	for job := range jobs {
		count.CountFilesInS3Folder(s3c, bucket, job)
	}
	defer wg.Done()
}
