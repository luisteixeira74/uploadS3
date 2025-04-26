package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	s3uploader "github.com/luisteixeira74/uploadS3/internal/s3uploader"
	"github.com/luisteixeira74/uploadS3/internal/uploader"
)

var (
	wg sync.WaitGroup
)

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		panic(err)
	}

	// Cria uma instância do S3Uploader
	s3 := &s3uploader.S3Uploader{
		Client: s3.New(sess),
		Bucket: os.Getenv("S3_BUCKET_NAME"),
	}

	dir, err := os.Open("./tmp")
	if err != nil {
		panic(err)
	}
	defer dir.Close()

	uploadControl := make(chan struct{}, 100) // 100 concurrent uploads
	errorFileUpload := make(chan string, 10)  // channel to handle error file uploads

	// goroutine to handle error file uploads
	go func() {
		for fullPath := range errorFileUpload {
			wg.Add(1)
			go func(path string) {
				defer wg.Done()
				uploadControl <- struct{}{}
				// Chama a função de upload com a interface
				go uploader.UploadFile(s3, path, uploadControl, errorFileUpload, &wg)
			}(fullPath)
		}
	}()

	// Leitura e upload dos arquivos
	for {
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

		fullPath := filepath.Join("./tmp", files[0].Name())
		go uploader.UploadFile(s3, fullPath, uploadControl, errorFileUpload, &wg)
	}

	// Fechar o canal e esperar pelas goroutines
	close(errorFileUpload)
	wg.Wait()
}
