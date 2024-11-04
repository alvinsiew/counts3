package worker

import (
	"sync"

	"github.com/alvinsiew/alvinsiew/counts3/internal/countfiles"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func worker(s3c *s3.Client, bucket string, jobs chan string, wg *sync.WaitGroup) {

	for job := range jobs {
		countFilesInS3Folder(s3c, bucket, job)
	}
	defer wg.Done()
}
