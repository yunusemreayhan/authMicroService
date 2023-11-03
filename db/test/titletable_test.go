package test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	dbutils "github.com/yunusemreayhan/goAuthMicroService/db"
	"github.com/yunusemreayhan/goAuthMicroService/db/sqlc"
)

func TestRoleDB(t *testing.T) {
	ctx := context.Background()

	con, err := pgx.Connect(ctx, dbutils.GetSQLDSN())
	if err != nil {
		t.Skip("db not connected !", err)
	}

	defer con.Close(ctx)

	queries := sqlc.New(con)
	_ = queries.DeleteTitleByName(ctx, "test")
	_, errDB := queries.CreateTitle(ctx, sqlc.CreateTitleParams{
		Name:        "test",
		Description: "test description",
	})

	if errDB != nil {
		t.Fatal(errDB)
	}

	_, errDB = queries.CreateTitle(ctx, sqlc.CreateTitleParams{
		Name:        "test",
		Description: "test description",
	})

	if errDB == nil {
		t.Fatal(errDB)
	}

	role, errDB := queries.GetTitleByName(ctx, "test")
	if errDB != nil {
		t.Fatal(errDB)
	}

	if role.Description != "test description" {
		t.Fatal("wrong role description " + role.Description)
	}

	// lazy clean up
	defer func() {
		errDB = queries.DeleteTitleByName(ctx, "test")
		if errDB != nil {
			t.Fatal(errDB)
		}
	}()
}
