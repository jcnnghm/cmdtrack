package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
)

// Code in this module is mostly copied from the golang examples

// Encrypt encrypts the provided data using AES-256-CBC, after generating a
// key by hashing the provided passphrase with sha256 and a salt
func Encrypt(data []byte, passphrase string) []byte {
	key := makeKeySlice(passphrase)

	// The scheme here appends the data with 1-16 bytes of padding, with every
	// byte being a uint8 of the padding length
	padLength := uint8(aes.BlockSize - (len(data) % aes.BlockSize))
	pad := make([]byte, padLength)
	for i := range pad {
		pad[i] = byte(padLength)
	}
	paddedData := append(data, pad...)

	// CBC mode works on blocks so plaintexts may need to be padded to the
	// next whole block. For an example of such padding, see
	// https://tools.ietf.org/html/rfc5246#section-6.2.3.2. Here we'll
	// assume that the plaintext is already of the correct length.
	if len(paddedData)%aes.BlockSize != 0 {
		panic("data is not a multiple of the block size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(paddedData))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], paddedData)

	return ciphertext
}

// Decrypt decrypts the provided ciphertext, reversing the Encrypt method.
func Decrypt(ciphertext []byte, passphrase string) []byte {
	key := makeKeySlice(passphrase)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	if len(ciphertext)%aes.BlockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(ciphertext, ciphertext)

	padLength := int(uint8(ciphertext[len(ciphertext)-1]))
	unpadded := ciphertext[:len(ciphertext)-padLength]

	return unpadded
}

func makeKeySlice(passphrase string) []byte {
	keyBytes := makeKey(passphrase)
	return keyBytes[:]
}

var salt = "cmdtrack!"

// Encrypt encrypts the provided data using AES-256-CBC, after generating a
func makeKey(passphrase string) [32]byte {
	return sha256.Sum256([]byte(salt + passphrase))
}
