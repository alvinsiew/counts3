package option

import (
	"flag"
	"fmt"
	"os"
)

// Flag allow for the cli
func Param() (string, string, int) {
	fileParam := flag.String("f", "", "File path and name")
	workerParam := flag.Int("w", 1, "Number of Workers")
	bucketParam := flag.String("s3", "", "S3 bucket name")
	flag.Parse()

	if *fileParam == "" {
		fmt.Println("File must be specific with -f")
		os.Exit(1)
	}

	if *bucketParam == "" {
		fmt.Println("S3 Bucket must be specific with -s3")
		os.Exit(1)
	}

	return *bucketParam, *fileParam, *workerParam
}
