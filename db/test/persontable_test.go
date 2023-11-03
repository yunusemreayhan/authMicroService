package test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	dbutils "github.com/yunusemreayhan/goAuthMicroService/db"
	"github.com/yunusemreayhan/goAuthMicroService/db/sqlc"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestPersonDB(t *testing.T) {
	ctx := context.Background()

	con, err := pgx.Connect(ctx, dbutils.GetSQLDSN())
	if err != nil {
		t.Skip("db not connected !", err)
	}

	defer con.Close(ctx)

	queries := sqlc.New(con)
	_, errDB := queries.CreatePerson(ctx, sqlc.CreatePersonParams{
		Personname:   "test",
		Email:        "test@test.com",
		PasswordHash: "testhash",
	})

	if errDB != nil {
		t.Fatal(errDB)
	}

	_, errDB = queries.CreatePerson(ctx, sqlc.CreatePersonParams{
		Personname:   "test",
		Email:        "test@test.com",
		PasswordHash: "testhash",
	})

	if errDB == nil {
		t.Fatal(errDB)
	}

	person, errDB := queries.GetPersonByEmail(ctx, "test@test.com")
	if errDB != nil {
		t.Fatal(errDB)
	}

	if person.PasswordHash != "testhash" {
		t.Fatal("wrong person password " + person.PasswordHash)
	}

	if person.Personname != "test" {
		t.Fatal("wrong person name " + person.Personname)
	}

	// lazy clean up
	defer func() {
		errDB = queries.DeletePersonByEmail(ctx, "test@test.com")
		if errDB != nil {
			t.Fatal(errDB)
		}
	}()
}
