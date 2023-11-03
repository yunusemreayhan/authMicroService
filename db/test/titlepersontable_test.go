package test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	dbutils "github.com/yunusemreayhan/goAuthMicroService/db"
	"github.com/yunusemreayhan/goAuthMicroService/db/sqlc"
)

func TestTitlePersonDB(t *testing.T) {
	ctx := context.Background()

	con, err := pgx.Connect(ctx, dbutils.GetSQLDSN())
	if err != nil {
		t.Skip("db not connected !", err)
	}

	defer con.Close(ctx)

	queries := sqlc.New(con)
	tempPerson, errDB := queries.CreatePerson(ctx, sqlc.CreatePersonParams{
		Personname:   "test",
		Email:        "test@test.com",
		PasswordHash: "testhash",
	})

	if errDB != nil {
		t.Fatal(errDB)
	}

	// lazy clean up
	defer func() {
		errDB := queries.DeletePersonByEmail(ctx, "test@test.com")

		if errDB != nil {
			t.Fatal(errDB)
		}
	}()

	tempTitle, errDB := queries.CreateTitle(ctx, sqlc.CreateTitleParams{
		Name:        "test",
		Description: "test description",
	})

	if errDB != nil {
		t.Fatal(errDB)
	}

	// lazy clean up
	defer func() {
		errDB := queries.DeleteTitleByName(ctx, "test")

		if errDB != nil {
			t.Fatal(errDB)
		}
	}()

	createdTitlePerson, errDB := queries.CreateTitlePerson(ctx, sqlc.CreateTitlePersonParams{
		PersonID: pgtype.Int4{Int32: tempPerson.ID, Valid: true},
		TitleID:  pgtype.Int4{Int32: tempTitle.ID, Valid: true},
	})

	if errDB != nil {
		t.Fatal(errDB)
	}

	if createdTitlePerson.PersonID.Int32 != tempPerson.ID {
		t.Fatal("wrong person id")
	}

	if createdTitlePerson.TitleID.Int32 != tempTitle.ID {
		t.Fatal("wrong title id")
	}

	// lazy clean up
	defer func() {
		errDB := queries.DeleteTitlePersonByPersonIDAndTitleID(ctx, sqlc.DeleteTitlePersonByPersonIDAndTitleIDParams{
			PersonID: pgtype.Int4{Int32: tempPerson.ID, Valid: true},
			TitleID:  pgtype.Int4{Int32: tempTitle.ID, Valid: true},
		})

		if errDB != nil {
			t.Fatal(errDB)
		}
	}()
}
