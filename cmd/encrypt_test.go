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

func ExampleEncrypt() {
	text := "Example Text to Encrypt"
	pass := "Passphrase!"
	encrypted := Encrypt([]byte(text), pass)
	fmt.Println(string(Decrypt(encrypted[:], pass)))
	// Output: Example Text to Encrypt
}
