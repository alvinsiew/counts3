# counts3
This cli help to count files in S3 folders concurrently using goroutines

# Usage

## Example

### To count number of files in each folders
counts3 -f <Filename + path> -w <Number of Workers to run concurrent job> -s3 <S3 Bucket Name>

### To get total count
counts3 -f <Filename + path> -w <Number of Workers to run concurrent job> -s3 <S3 Bucket Name> | awk '{print $6}' | xargs  | sed -e 's/\ /+/g' | bc

## MacOS
$ env GOOS=darwin GOARCH=amd64 go build -o counts3 main.go
