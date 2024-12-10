package main

import (
	"crypto/sha512"
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: artifact-key-generator <filename>")
		os.Exit(1)
	}

	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	hash := sha512.New()
	if _, err := io.Copy(hash, file); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("Error getting file info: %v\n", err)
		os.Exit(1)
	}
	hashString := fmt.Sprintf("%x", hash.Sum(nil))
	fileSize := fileInfo.Size()
	fmt.Printf("%s-%d\n", hashString, fileSize)
}