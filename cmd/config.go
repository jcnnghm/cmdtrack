package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"strings"
)

var Config configStruct

type configStruct struct {
	SharedSecret  string `json:"shared-secret"`
	EncryptionKey string `json:"encryption-key"`
}

func LoadConfig() {
	usr, _ := user.Current()
	dir := usr.HomeDir
	file := dir + "/.cmdtrack.conf"

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		panic(file + " must exist")
	}

	fileInfo, err := os.Stat(file)
	if err != nil {
		panic(file + " must exist")
	}
	if fileInfo.Mode().Perm() != 0600 {
		panic("Mode of " + file + " must be 600")
	}

	dec := json.NewDecoder(strings.NewReader(string(bytes)))
	if err := dec.Decode(&Config); err != nil {
		panic("Parsing ~/.cmdtrack.conf failed")
	}
}
