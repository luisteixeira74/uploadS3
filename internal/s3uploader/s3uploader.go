package s3uploader

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Uploader implementa a interface Uploader
type S3Uploader struct {
	Client *s3.S3
	Bucket string
}

// Implementa a interface Uploader
func (u *S3Uploader) Upload(filename string, body io.ReadSeeker) error {
	_, err := u.Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(u.Bucket),
		Key:    aws.String(filename),
		Body:   body,
	})
	return err
}
