package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func createFiles(numFiles int, dir string) error {
	// Cria o diretório caso não exista
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < numFiles; i++ {
		filename := fmt.Sprintf("%s/file_%d_%d.txt", dir, time.Now().UnixNano(), rng.Intn(1000))
		file, err := os.Create(filename)
		if err != nil {
			fmt.Printf("Erro criando arquivo: %s\n", err)
			continue
		}
		file.Close()
		fmt.Println("Criado:", filename)
	}
	return nil
}

func main() {
	const numFiles = 50 // Quantidade de arquivos a criar
	err := createFiles(numFiles, "./tmp")
	if err != nil {
		panic(err)
	}
}
