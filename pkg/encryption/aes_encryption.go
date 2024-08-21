package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"strings"
)

type AesEncryption struct {
}

func NewAesEncryption() *AesEncryption {
	return &AesEncryption{}
}

func (c *AesEncryption) Encrypt(plaintext string) (string, error) {
	if strings.TrimSpace(plaintext) == "" {
		return plaintext, nil
	}

	key, err := getKeyFromEnv()
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return hex.EncodeToString(ciphertext), nil
}

func (c *AesEncryption) Decrypt(ciphertext string) (string, error) {
	if strings.TrimSpace(ciphertext) == "" {
		return ciphertext, nil
	}

	key, err := getKeyFromEnv()
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertextBytes, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	if len(ciphertextBytes) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	iv := ciphertextBytes[:aes.BlockSize]
	ciphertextBytes = ciphertextBytes[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertextBytes, ciphertextBytes)

	return string(ciphertextBytes), nil
}

func (c *AesEncryption) EncryptMultiple(fields map[string]string) (map[string]string, error) {
	encryptedFields := make(map[string]string)
	for key, value := range fields {
		encryptedValue, err := c.Encrypt(value)
		if err != nil {
			return nil, err
		}
		encryptedFields[key] = encryptedValue
	}
	return encryptedFields, nil
}

func getKeyFromEnv() ([]byte, error) {
	key := os.Getenv("AES_KEY")
	if key == "" {
		return nil, errors.New("AES_KEY environment variable is not set")
	}
	return hex.DecodeString(key)
}
