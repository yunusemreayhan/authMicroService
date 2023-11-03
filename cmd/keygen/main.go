package main

import (
	"crypto/rsa"
	"flag"
	"fmt"

	"github.com/yunusemreayhan/goAuthMicroService/internal/key"
)

func main() {
	keyPath := flag.String("key", "./private.key", "path to private key")
	var newkey *rsa.PrivateKey
	flag.Parse()
	fmt.Printf("generating new private key with keyPath: %v\n", *keyPath)
	key.StorePrivateKey(newkey, *keyPath)
}
