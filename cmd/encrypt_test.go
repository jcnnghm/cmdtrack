package cmd

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	cases := []string{"exampleplaintext", "another", "t"}
	pass := "example"
	for _, testCase := range cases {
		data := []byte(testCase)
		encrypted := Encrypt(data, pass)

		assert.False(t, bytes.Equal(data, encrypted),
			"Encrypted value should be different")
		assert.Equal(t, data, Decrypt(encrypted[:], pass),
			"%v did not encrypt cleanly.", testCase)
	}
}

func TestEncryptBase64(t *testing.T) {
	cases := []string{"exampleplaintext", "another", "t"}
	pass := "example"
	for _, testCase := range cases {
		encrypted := EncryptBase64(testCase, pass)

		assert.NotEqual(t, testCase, encrypted,
			"Encrypted value should be different")

		decrypted, err := DecryptBase64(encrypted, pass)
		assert.Nil(t, err)
		assert.Equal(t, testCase, decrypted,
			"%v did not encrypt cleanly.", testCase)
	}
}

func ExampleEncrypt() {
	text := "Example Text to Encrypt"
	pass := "Passphrase!"
	encrypted := Encrypt([]byte(text), pass)
	fmt.Println(string(Decrypt(encrypted[:], pass)))
	// Output: Example Text to Encrypt
}
