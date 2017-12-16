package main

// TODO: none of this is probably needed. jwt-go will probably replace
// it if it hasnt already

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"io/ioutil"
	"os"
)

type Token struct {
	filename   string
	privateKey *rsa.PrivateKey
}

func (t *Token) Sign(data []byte) ([]byte, error) {
	return t.privateKey.Sign(rand.Reader, data, crypto.SHA256)
}

func NewToken(filename string) (*Token, error) {
	pKey, err := GetPrivateKey(filename)
	if err != nil {
		return &Token{}, err
	}

	return &Token{
		filename:   filename,
		privateKey: pKey,
	}, nil
}

func GetPrivateKey(filename string) (*rsa.PrivateKey, error) {
	if filename == "" {
		filename = ".local/bitwarden-go/jwt-rsa.key"
	}

	// Read existing file
	keyData, err := ioutil.ReadFile(filename)
	if err == nil {
		return x509.ParsePKCS1PrivateKey(keyData)
	}

	if os.IsNotExist(err) {
		return CreatePrivateKey(filename)
	}

	return nil, err
}

func CreatePrivateKey(filename string) (*rsa.PrivateKey, error) {
	// File does exist, create and return
	pKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	return pKey, WritePrivateKey(filename, pKey)
}

func WritePrivateKey(filename string, pKey *rsa.PrivateKey) error {
	keyData := x509.MarshalPKCS1PrivateKey(pKey)
	return ioutil.WriteFile(filename, keyData, 0600)
}
