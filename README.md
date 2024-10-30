# counts3
This cli help to count files in S3 folders concurrently using goroutines

# Example
counts3 -f <Filename + path> -w <Number of Workers to run concurrent job> -s3 <S3 Bucket Name>

counts3 -f directorylist.txt -w 3 -s3 S3Name