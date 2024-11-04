package count

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

// countFilesInS3Folder counts the number of objects in a specified S3 folder
func CountFilesInS3Folder(client *s3.Client, bucket string, prefix string) {
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
