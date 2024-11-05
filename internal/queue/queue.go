package queue

import (
	"bufio"
	"log"
	"os"
)

// Read file line by line and send to channel
func QueueJob(fileName string, jobs chan string) {
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
}
