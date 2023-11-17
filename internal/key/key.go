package key

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"log"
	"os"
)

const DefaultKeyPath = "/private.key"

var PrivateKey *rsa.PrivateKey

func GeneratePrivateKey() *rsa.PrivateKey {
	rng := rand.Reader
	privateKey, err := rsa.GenerateKey(rng, 2048)
	if err != nil {
		log.Fatalf("rsa.GenerateKey: %v", err)
	}
	return privateKey
}

func StorePrivateKey(privateKey *rsa.PrivateKey, path_for_private_key string) {
	if privateKey == nil {
		privateKey = GeneratePrivateKey()
	}
	// store private key
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	err := os.WriteFile(path_for_private_key, privateKeyBytes, 0600)
	if err != nil {
		log.Fatalf("os.WriteFile: %v", err)
	}

	// store public key
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		log.Fatalf("x509.MarshalPKIXPublicKey: %v", err)
	}

	err = os.WriteFile(path_for_private_key+".pub", publicKeyBytes, 0600)
	if err != nil {
		log.Fatalf("os.WriteFile: %v", err)
	}
}

func LoadPrivateKey(path_for_private_key string) (*rsa.PrivateKey, error) {
	// load private key
	privateKeyBytes, err := os.ReadFile(path_for_private_key)
	if err != nil {
		return nil, err
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBytes)
	if err != nil {
		return nil, err
	}

	log.Default().Printf("returning loaded privateKey")

	return privateKey, nil
}
