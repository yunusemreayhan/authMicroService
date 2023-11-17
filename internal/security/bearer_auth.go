package security

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	dbutils "github.com/yunusemreayhan/goAuthMicroService/db"
	"github.com/yunusemreayhan/goAuthMicroService/db/sqlc"
	"github.com/yunusemreayhan/goAuthMicroService/internal/config"
	"github.com/yunusemreayhan/goAuthMicroService/internal/key"
	"github.com/yunusemreayhan/goAuthMicroService/internal/model"
)

func ApiKeyAuthHandler(Token string) (interface{}, error) {
	var ret sqlc.Person
	log.Default().Printf("ApiKeyAuthHandler")

	privateKey, err := key.LoadPrivateKey(config.GetConfig().DefaultKeyPath)
	if err != nil {
		log.Default().Printf("LoadPrivateKey error: [%v]\n", err)
		return nil, fmt.Errorf("LoadPrivateKey error: [%v]", err)
	}

	token, err := jwt.Parse(Token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			log.Default().Printf("unexpected signing method: [%v]", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: [%v]", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return privateKey.Public(), nil
	})

	if err == nil {
		if token == nil {
			log.Default().Printf("ApiKeyAuthHandler token is nil\n")
			return nil, fmt.Errorf("invalid token")
		}
		var identity model.JWTIdentity
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			identity.Id = int32(claims["identifier"].(map[string]interface{})["id"].(float64))
			identity.Username = claims["identifier"].(map[string]interface{})["username"].(string)
			identity.Email = claims["identifier"].(map[string]interface{})["email"].(string)
		} else {
			log.Default().Printf("ApiKeyAuthHandler claims error : [%v]\n", err)
			return nil, errors.New("unauthorized")
		}

		ctx := context.Background()

		con, err := pgx.Connect(ctx, dbutils.GetSQLDSN())
		if err != nil {
			log.Default().Printf("ApiKeyAuthHandler pgx.Connect error : [%v]\n", err)
			return nil, errors.New("unauthorized")
		}

		defer func(con *pgx.Conn, ctx context.Context) {
			err := con.Close(ctx)
			if err != nil {
				log.Default().Printf("ApiKeyAuthHandler con.Close error : [%v]\n", err)
			}
		}(con, ctx)

		queries := sqlc.New(con)

		ret, err = queries.GetPerson(ctx, identity.Id)
		if err == nil {
			return &ret, nil
		}
	}

	log.Default().Printf("ApiKeyAuthHandler err : [%v]\n", err)
	return nil, err
}
