package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

func encryptAES(data, key, iv []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	ivCopy := make([]byte, aes.BlockSize)
	copy(ivCopy, iv)
	stream := cipher.NewCFBEncrypter(block, ivCopy)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return base64.StdEncoding.EncodeToString(append(ivCopy, ciphertext[aes.BlockSize:]...)), nil
}
