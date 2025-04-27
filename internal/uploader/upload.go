package uploader

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Uploader é a interface que define o comportamento de upload de arquivos.
type Uploader interface {
	Upload(filename string, body io.ReadSeeker) error
}

func UploadFile(uploader Uploader, fullPath string, uploadControl chan struct{}, errorFileUpload chan<- string) error {
	defer func() {
		if uploadControl != nil {
			select {
			case <-uploadControl:
				// Liberou slot
			default:
			}
		}
	}()

	filename := filepath.Base(fullPath)
	fmt.Printf("Uploading file %s started\n", filename)

	const maxRetries = 2 // número máximo de tentativas (retry c/ backoff exponencial)

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		file, err := os.Open(fullPath)
		if err != nil {
			lastErr = err
			fmt.Printf("Erro ao abrir o arquivo %s: %v\n", fullPath, err)
			break // não adianta tentar de novo se nem conseguiu abrir
		}
		defer file.Close()

		err = uploader.Upload(filename, file)

		if err == nil {
			fmt.Printf("Successfully uploaded %s\n", filename)
			return nil
		}

		fmt.Printf("Erro no upload (%s), tentativa %d de %d: %v\n", filename, attempt+1, maxRetries, err)
		lastErr = err

		// Backoff exponencial
		sleepDuration := time.Duration(1<<attempt) * time.Second
		time.Sleep(sleepDuration)
	}

	fmt.Printf("Falha no upload de %s após %d tentativas: %v\n", filename, maxRetries, lastErr)
	errorFileUpload <- fullPath // envia para reprocessamento
	return lastErr              // Retorna o erro final para que o código chamador possa tratá-lo
}
