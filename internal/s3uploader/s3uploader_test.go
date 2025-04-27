package s3uploader

import (
	"io"
	"os"
	"testing"

	"github.com/luisteixeira74/uploadS3/internal/uploader"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TestUploadFileSuite struct {
	suite.Suite
	mockUploader *MockUploader
}

type MockUploader struct {
	mock.Mock
}

func (m *MockUploader) Upload(filename string, body io.ReadSeeker) error {
	args := m.Called(filename, body)
	return args.Error(0)
}

func (suite *TestUploadFileSuite) SetupTest() {
	// Setup comum para todos os testes (se necessário)
	suite.mockUploader = new(MockUploader)
}

func (suite *TestUploadFileSuite) TestUploadFileSuccess() {
	suite.mockUploader.On("Upload", "test.txt", mock.Anything).Return(nil)

	// Criação de um arquivo de teste
	file, err := os.Create("/tmp/test.txt")
	suite.Require().NoError(err)

	defer file.Close()
	file.WriteString("Test file content")

	// Testando a função de upload
	uploadControl := make(chan struct{}, 1)
	errorFileUpload := make(chan string, 1)

	uploadControl <- struct{}{}

	// Chama a função de upload
	err = uploader.UploadFile(suite.mockUploader, "/tmp/test.txt", uploadControl, errorFileUpload)
	suite.NoError(err)

	// Asserções
	suite.mockUploader.AssertExpectations(suite.T())
	select {
	case errFile := <-errorFileUpload:
		suite.Fail("Não esperava erro de upload, mas recebeu: %s", errFile)
	default:
		// Nenhum erro
	}
}

func (suite *TestUploadFileSuite) TestUploadFileWithRetry() {
	// Similar ao primeiro teste, mas com lógica de retry
}

func TestRunUploadFileSuite(t *testing.T) {
	suite.Run(t, new(TestUploadFileSuite))
}
