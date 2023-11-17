package security

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	dbutils "github.com/yunusemreayhan/goAuthMicroService/db"
	"github.com/yunusemreayhan/goAuthMicroService/db/sqlc"
	"log"
)

func BasicAuthHandler(username string, password string) (interface{}, error) {
	ctx := context.Background()

	con, err := pgx.Connect(ctx, dbutils.GetSQLDSN())
	if err != nil {
		return nil, errors.New("unauthorized")
	}

	defer func(con *pgx.Conn, ctx context.Context) {
		err := con.Close(ctx)
		if err != nil {
			log.Default().Printf("con.Close error : [%v]\n", err)
		}
	}(con, ctx)

	queries := sqlc.New(con)

	voucherOwnerPerson, err := queries.GetPersonByPersonNameAndPassword(ctx, sqlc.GetPersonByPersonNameAndPasswordParams{
		Personname:   username,
		PasswordHash: password, // todo query with password hash
	})
	if err != nil {
		return nil, errors.New("unauthorized")
	}
	return &voucherOwnerPerson, nil
}
