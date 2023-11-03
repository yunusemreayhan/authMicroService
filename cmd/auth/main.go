package main

import (
	"log"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/yunusemreayhan/goAuthMicroService/internal/handler"
	"github.com/yunusemreayhan/goAuthMicroService/internal/key"
)

func main() {
	app := fiber.New()

	// load private key
	privateKey, err := key.LoadPrivateKey(key.DefaultKeyPath)
	if err != nil {
		log.Fatalf("key.LoadPrivateKey: %v", err)
	}

	// Custom config
	app.Static("/", "../../public", fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        true,
		Index:         "index.html",
		CacheDuration: -1,
		MaxAge:        0,
	})

	// Authentication Microservice
	app.Post("/api/register", handler.RegisterPerson)
	app.Post("/api/login", handler.LoginPerson)
	app.Post("/api/verify", handler.VerifyToken)

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.RS256,
			Key:    privateKey.Public(),
		},
	}))

	// Person Database Microservice
	app.Get("/api/person", handler.UpdatePerson)
	app.Post("/api/person", handler.UpdatePerson)

	app.Listen(":3000")
}
