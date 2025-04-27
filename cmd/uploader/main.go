package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
	s3uploader "github.com/luisteixeira74/uploadS3/internal/s3uploader"
	"github.com/luisteixeira74/uploadS3/internal/uploader"
)

var (
	wg sync.WaitGroup
)

func main() {
	// Carregar as variáveis de ambiente do arquivo .env
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Erro ao carregar o arquivo .env")
		return
	}

	// Obter as credenciais da AWS a partir das variáveis de ambiente
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")

	if accessKey == "" || secretKey == "" || region == "" {
		fmt.Println("Erro: Credenciais AWS ou região não configuradas corretamente.")
		return
	}

	// Criar a sessão da AWS com as credenciais e região
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})
	if err != nil {
		fmt.Printf("Erro ao criar sessão AWS: %v\n", err)
		return
	}

	// Cria uma instância do S3Uploader
	s3 := &s3uploader.S3Uploader{
		Client: s3.New(sess),
		Bucket: os.Getenv("S3_BUCKET_NAME"),
	}

	// Teste para ver se conseguimos acessar o S3
	fmt.Println("Sessão AWS criada com sucesso!")

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
			uploadControl <- struct{}{}
			go func(path string) {
				defer wg.Done()
				// Chama a função de upload com a interface
				err := uploader.UploadFile(s3, path, uploadControl, errorFileUpload)
				if err != nil {
					fmt.Printf("Falha ao tentar fazer o upload: %v\n", err)
					// Lidar com o erro ou tentar novamente, etc.
				}
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

		go func(fullPath string) {
			defer wg.Done()
			// Chama a função de upload com a interface
			err := uploader.UploadFile(s3, fullPath, uploadControl, errorFileUpload)
			if err != nil {
				fmt.Printf("Falha ao tentar fazer o upload: %v\n", err)
				// Lidar com o erro ou tentar novamente, etc.
			}
		}(fullPath)
	}

	// Fechar o canal e esperar pelas goroutines
	go func() {
		wg.Wait()
		close(errorFileUpload)
	}()

	// Mensagem informando que o processo foi concluído
	fmt.Println("Processo de upload concluído.")
}
