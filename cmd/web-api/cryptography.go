package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

func (app *application) encryptSecret(textToEncrypt string) ([]byte,error) {
	plainText := []byte(textToEncrypt)
	
	block, err := aes.NewCipher(app.cryptoKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	cipherText := gcm.Seal(nonce, nonce, plainText, nil)
	fmt.Println("The ciphered text: " + string(cipherText))
	
	return cipherText, nil
}

func (app *application) decryptSecret(secretToDecrypt []byte) (string, error) {
	block, err := aes.NewCipher(app.cryptoKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := (secretToDecrypt)[:gcm.NonceSize()]
	cipherText := (secretToDecrypt)[gcm.NonceSize():]

	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}