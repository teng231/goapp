package crypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// PKCS7 padding
func pkcs7Pad(b []byte, blockSize int) []byte {
	padLen := blockSize - (len(b) % blockSize)
	if padLen == 0 {
		padLen = blockSize
	}
	pad := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(b, pad...)
}

func pkcs7Unpad(b []byte, blockSize int) ([]byte, error) {
	if len(b) == 0 || len(b)%blockSize != 0 {
		return nil, errors.New("invalid padded data")
	}
	padLen := int(b[len(b)-1])
	if padLen == 0 || padLen > blockSize {
		return nil, errors.New("invalid padding size")
	}
	// verify padding bytes
	for i := 0; i < padLen; i++ {
		if b[len(b)-1-i] != byte(padLen) {
			return nil, errors.New("invalid padding")
		}
	}
	return b[:len(b)-padLen], nil
}

// Encrypt plaintext using AES-CBC with PKCS7 padding.
// key must be 16, 24 or 32 bytes long (AES-128/192/256).
// Returns base64( IV || ciphertext ).
func EncryptAES_CBC_Base64(plaintext []byte, key []byte) (string, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return "", errors.New("key length must be 16, 24 or 32 bytes")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize() // 16

	// PKCS7 pad
	padded := pkcs7Pad(plaintext, blockSize)

	// random IV
	iv := make([]byte, blockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	ciphertext := make([]byte, len(padded))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, padded)

	// Prepend IV to ciphertext, then base64 encode
	out := append(iv, ciphertext...)
	return base64.StdEncoding.EncodeToString(out), nil
}

// Decrypt base64( IV || ciphertext ) using AES-CBC with PKCS7
func DecryptAES_CBC_Base64(b64 string, key []byte) ([]byte, error) {
	raw, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, errors.New("key length must be 16, 24 or 32 bytes")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize() // 16

	if len(raw) < blockSize || len(raw)%blockSize != 0 {
		return nil, errors.New("ciphertext too short or not multiple of blocksize")
	}

	iv := raw[:blockSize]
	ciphertext := raw[blockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	plaintextPadded := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintextPadded, ciphertext)

	// unpad
	return pkcs7Unpad(plaintextPadded, blockSize)
}
