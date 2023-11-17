package main

import (
	"github.com/yunusemreayhan/goAuthMicroService/internal/handler"
	"github.com/yunusemreayhan/goAuthMicroService/internal/key"
	"log"
)

func main() {
	var restApp handler.RestAPPForAuth
	restApp.PrepareFiberApp()
	key, err := key.LoadPrivateKey(key.DefaultKeyPath)
	if err != nil {
		log.Fatal("unable to load private key", err)
	}
	restApp.SetPrivateKey(key)
	restApp.GetFiberApp().Listen(":3000")
}
