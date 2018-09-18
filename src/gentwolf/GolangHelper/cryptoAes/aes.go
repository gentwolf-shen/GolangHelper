package cryptoAes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

func Encrypt(src []byte, key, iv string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	src = PKCSPadding(src, block.BlockSize())
	data := make([]byte, len(src))

	encrypter := cipher.NewCBCEncrypter(block, []byte(iv))
	encrypter.CryptBlocks(data, src)

	return base64.StdEncoding.EncodeToString(data), nil
}

func Decrypt(src, key, iv string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	decrypter := cipher.NewCBCDecrypter(block, []byte(iv))

	bytes := base64Decode(src)
	data := make([]byte, len(bytes))
	decrypter.CryptBlocks(data, bytes)

	return PKCSUnPadding(data), nil
}

func base64Decode(str string) []byte {
	b, _ := base64.StdEncoding.DecodeString(str)
	return b
}

func PKCSPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCSUnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}
