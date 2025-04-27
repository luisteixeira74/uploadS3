package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateFiles(t *testing.T) {
	// Testa a criação de arquivos com sucesso

	// Cria o diretório dentro de uma pasta "test" dedicada
	dir := "../../test/tmp_test"

	// Garante que o diretório será limpo após o teste
	defer os.RemoveAll(dir)

	// Chama a função que cria os arquivos
	err := createFiles(5, dir) // Cria 5 arquivos de teste

	// Verifica se não houve erro na criação
	require.NoError(t, err)

	// Verifica se o diretório foi criado
	info, err := os.Stat(dir)
	require.NoError(t, err)
	assert.True(t, info.IsDir(), "O diretório deve ser criado")

	// Verifica se os arquivos foram criados
	files, err := os.ReadDir(dir)
	require.NoError(t, err)
	assert.Len(t, files, 5, "Deve criar exatamente 5 arquivos")
}
