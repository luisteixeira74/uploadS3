package main

import (
	"fmt"
	"os"
)

func main() {
	i := 0
	for {
		f, err := os.Create(fmt.Sprintf("./tmp/file%d.txt", i))
		if err != nil {
			fmt.Println("Error creating file:", err)
			break
		}
		defer f.Close()
		f.WriteString(fmt.Sprintf("This is file number %d\n", i))
		i++
	}
}
