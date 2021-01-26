package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const ext = ".c8"

func decrypt(cipherstring []byte, keystring string) string {
	ciphertext := cipherstring

	// Key
	key := []byte(keystring)

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// Before even testing the decryption,
	// if the text is too small, then it is incorrect
	if len(ciphertext) < aes.BlockSize {
		panic("Text is too short")
	}

	// Get the 16 byte IV
	iv := ciphertext[:aes.BlockSize]

	// Remove the IV from the ciphertext
	ciphertext = ciphertext[aes.BlockSize:]

	// Return a decrypted stream
	stream := cipher.NewCFBDecrypter(block, iv)

	// Decrypt bytes from ciphertext
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext)
}

func encrypt(data []byte, keystring string) string {
	// Key
	key := []byte(keystring)

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// Empty array of 16 + plaintext length
	// Include the IV at the beginning
	ciphertext := make([]byte, aes.BlockSize+len(data))

	// Slice of first 16 bytes
	iv := ciphertext[:aes.BlockSize]

	// Write 16 rand bytes to fill iv
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	// Return an encrypted stream
	stream := cipher.NewCFBEncrypter(block, iv)

	// Encrypt bytes from plaintext to ciphertext
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return string(ciphertext)
}

func writeToFile(data, file string) {
	ioutil.WriteFile(file, []byte(data), 777)
}

func readFromFile(file string) ([]byte, error) {
	data, err := ioutil.ReadFile(file)
	return data, err
}

func encryptFiles(dir string, key string) {
	// get files in folder
	files, err := filepath.Glob(dir)
	if err != nil {
		log.Fatal(err)
	}

	// for each file
	// encrypt, write new file, delete old file
	for _, filePath := range files {
		file, err := os.Lstat(filePath)
		if err != nil {
			log.Fatal(err)
		}

		if file.IsDir() {
			continue
		}

		if file.Name()[len(file.Name())-3:] == ext {
			continue
		}

		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatal(err)
		}

		cipherData := encrypt(data, key)

		newFileName := filePath + ext
		writeToFile(cipherData, newFileName)

		_ = os.Remove(filePath)

		fmt.Println(file.Name())
	}
}

func decryptFiles(dir string, key string) {
	// get files in folder
	files, err := filepath.Glob(dir)
	if err != nil {
		log.Fatal(err)
	}

	// for each file
	// decrypt, write new file, delete old file
	for _, filePath := range files {
		file, err := os.Lstat(filePath)
		if err != nil {
			log.Fatal(err)
		}

		if file.IsDir() {
			continue
		}

		if file.Name()[len(file.Name())-3:] != ext {
			continue
		}

		data, err := readFromFile(filePath)
		if err != nil {
			log.Fatal(err)
		}

		decryptedData := decrypt(data, key)

		newFileName := strings.ReplaceAll(filePath, ext, "")
		writeToFile(decryptedData, newFileName)

		_ = os.Remove(filePath)

		fmt.Println(file.Name())
	}
}

func main() {
	cmd := os.Args[1]
	dir := os.Args[2]
	key := os.Args[3]

	switch cmd {
	case "encrypt":
		encryptFiles(dir, key)
		break

	case "decrypt":
		decryptFiles(dir, key)
		break
	}
}
