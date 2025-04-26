package s3uploader

import "io"

// Uploader é a interface que define o comportamento de upload de arquivos.
type Uploader interface {
	Upload(filename string, body io.ReadSeeker) error
}
