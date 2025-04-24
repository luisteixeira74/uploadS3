package main

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	s3Client *s3.S3
	s3Bucket string
	wg       sync.WaitGroup
)

func init() {
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String("us-east-1"),
			Credentials: credentials.NewStaticCredentials(
				os.Getenv("AWS_ACCESS_KEY_ID"),
				os.Getenv("AWS_SECRET_ACCESS_KEY"),
				"",
			),
		},
	)
	if err != nil {
		panic(err)
	}
	// Create S3 client
	s3Client = s3.New(sess)
	s3Bucket = os.Getenv("S3_BUCKET_NAME")
}

func main() {
	dir, err := os.Open("./tmp")
	if err != nil {
		panic(err)
	}
	defer dir.Close()

	uploadControl := make(chan struct{}, 100) // 100 concurrent uploads
	errorFileUpload := make(chan string, 10)  // channel to handle error file uploads

	// goroutine to handle error file uploads
	// This goroutine will listen for error messages and retry the upload
	// when a file fails to upload
	go func() {
		for filename := range errorFileUpload {
			wg.Add(1)
			go func(fname string) {
				defer wg.Done()
				uploadControl <- struct{}{}
				uploadFile(fname, uploadControl, errorFileUpload)
			}(filename)
		}
	}()

	for {
		// ReadDir returns a slice of FileInfo
		files, err := dir.ReadDir(1)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Error reading directory: %s\n", err)
			continue
		}

		wg.Add(1)
		uploadControl <- struct{}{}
		uploadFile(files[0].Name(), uploadControl, errorFileUpload)
	}
	wg.Wait()
}

func uploadFile(filename string, uploadControl chan struct{}, errorfileUpload chan<- string) {
	defer wg.Done()
	completeFileName := fmt.Sprintf("./tmp/%s", filename)
	fmt.Printf("Uploading file %s started\n", filename)
	file, err := os.Open(completeFileName)
	if err != nil {
		fmt.Printf("Error opening file %s\n", filename)
		<-uploadControl                     // empty the channel
		errorfileUpload <- completeFileName // send the filename to the error channel
		return
	}
	defer file.Close()
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(filename),
		Body:   file,
	})
	if err != nil {
		fmt.Printf("Error uploading file %s to S3: %s\n", filename, err)
		<-uploadControl                     // empty the channel
		errorfileUpload <- completeFileName // send the filename to the error channel
		return
	}
	fmt.Printf("Successfully uploaded %s \n", filename)
	<-uploadControl // empty the channel
}
