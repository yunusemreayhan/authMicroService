package key

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"log"
	"os"
)

var PrivateKey *rsa.PrivateKey

func generatePrivateKey() *rsa.PrivateKey {
	rng := rand.Reader
	privateKey, err := rsa.GenerateKey(rng, 2048)
	if err != nil {
		log.Fatalf("rsa.GenerateKey: %v", err)
	}
	return privateKey
}

func StorePrivateKey(privateKey *rsa.PrivateKey, path_for_private_key string) {
	if privateKey == nil {
		privateKey = generatePrivateKey()
	}
	// store private key
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	err := os.WriteFile(path_for_private_key, privateKeyBytes, 0600)
	if err != nil {
		log.Printf("os.WriteFile: %v", err)
	}

	// store public key
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		log.Printf("x509.MarshalPKIXPublicKey: %v", err)
	}

	err = os.WriteFile(path_for_private_key+".pub", publicKeyBytes, 0600)
	if err != nil {
		log.Printf("os.WriteFile: %v", err)
	}
}

func LoadPrivateKey(pathForPrivateKey string) (*rsa.PrivateKey, error) {
	// load private key
	if PrivateKey != nil {
		return PrivateKey, nil
	}

	log.Default().Printf("loading privateKey")
	privateKeyBytes, err := os.ReadFile(pathForPrivateKey)
	if err != nil {
		log.Default().Printf("GENERATE PRIVATE KEY")
		key := generatePrivateKey()
		StorePrivateKey(key, pathForPrivateKey)
		PrivateKey = key
	} else {
		privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBytes)
		if err != nil {
			log.Fatalf("x509.ParsePKCS1PrivateKey: %v", err)
		}
		PrivateKey = privateKey
	}

	log.Default().Printf("returning loaded privateKey")

	return PrivateKey, nil
}
