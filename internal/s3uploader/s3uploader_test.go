package s3uploader

import (
	"io"
	"sync"
	"testing"

	"github.com/luisteixeira74/uploadS3/internal/uploader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var wg sync.WaitGroup

// MockUploader simula a interface Uploader
type MockUploader struct {
	mock.Mock
}

// Alteração do tipo do parâmetro de io.Reader para io.ReadSeeker
func (m *MockUploader) Upload(filename string, body io.ReadSeeker) error {
	args := m.Called(filename, body)
	return args.Error(0)
}

// TestUploadFileSuccess testa um upload bem-sucedido
func TestUploadFileSuccess(t *testing.T) {
	m := new(MockUploader)
	m.On("Upload", "test.txt", mock.Anything).Return(nil)

	uploadControl := make(chan struct{}, 1)
	errorFileUpload := make(chan string, 1)

	uploadControl <- struct{}{} // simula slot ocupado

	// Prepara o WaitGroup para aguardar uma goroutine
	wg.Add(1)

	// Lança a função assíncrona para testar
	go func() {
		defer wg.Done() // Garante que o WaitGroup será decrementado após a execução
		// Chamando a função que você quer testar
		uploader.UploadFile(m, "testdata/test.txt", uploadControl, errorFileUpload, nil)
	}()

	// Aguarda a execução das goroutines antes de verificar as asserções
	wg.Wait()

	// Asserções
	m.AssertExpectations(t)
	assert.Empty(t, errorFileUpload) // Verifica se o canal de erro está vazio (nenhum erro ocorreu)
}
