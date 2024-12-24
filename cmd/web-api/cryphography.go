package main

// import (
// 	"crypto/aes"
// 	"crypto/cipher"
// 	"crypto/rand"
// 	"fmt"
// )

// func (app *application) encryptSecret(plainText string) ([]byte, error) {
// 	block, err := aes.NewCipher(app.encryptionKey)
// 	if err != nil {
// 		return nil, err
// 	}

// 	gcm, err := cipher.NewGCM(block)
// 	if err != nil {
// 		return nil, err
// 	}

// 	nonce := make([]byte, gcm.NonceSize())
// 	_, err = rand.Read(nonce)
// 	if err != nil {
// 		return nil, err
// 	}

// 	cipherText := gcm.Seal(nonce, nonce, []byte(plainText), nil)
// 	fmt.Println("Cipher text:", cipherText)

// 	return cipherText, nil
// }
