package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	const numFiles = 20 // Quantidade de arquivos a criar

	// Cria o diretório caso não exista, com permissões completas (equivalente a chmod 777)
	err := os.MkdirAll("./tmp", os.ModePerm)
	if err != nil {
		panic(err)
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < numFiles; i++ {
		filename := fmt.Sprintf("./tmp/file_%d_%d.txt", time.Now().UnixNano(), rng.Intn(1000))
		file, err := os.Create(filename)
		if err != nil {
			fmt.Printf("Erro criando arquivo: %s\n", err)
			continue
		}
		file.Close()
		fmt.Println("Criado:", filename)
	}
}
