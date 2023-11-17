package auth_handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/jackc/pgx/v5"
	dbutils "github.com/yunusemreayhan/goAuthMicroService/db"
	"github.com/yunusemreayhan/goAuthMicroService/db/sqlc"
	"github.com/yunusemreayhan/goAuthMicroService/models"
	"github.com/yunusemreayhan/goAuthMicroService/restapi/operations/auth"
)

// curl -X POST -H "Content-Type: application/json" -d '{"person_name":"test4", "email":"test4@test", "password":"test4"}' localhost:3000/api/register
// curl -X POST -H "Content-Type: application/json" -d '{"person_name":"test4", "email":"test4@test", "password":"test4"}' localhost/api/register

type RegisterHandler struct {
	middleware.Responder
	req auth.RegisterParams
}

func (l *RegisterHandler) BindRequest(req auth.RegisterParams) {
	l.req = req
}

func (l *RegisterHandler) WriteResponse(rw http.ResponseWriter, pr runtime.Producer) {
	ctx := context.Background()

	con, err := pgx.Connect(ctx, dbutils.GetSQLDSN())
	if err != nil {
		log.Default().Printf("pgx.Connect error : [%v]\n", err)
		errStr := fmt.Sprintf("pgx.Connect error : [%v]", err)
		rw.WriteHeader(http.StatusInternalServerError)
		var ret models.Error
		ret.Code = http.StatusInternalServerError
		ret.Message = &errStr
		json.NewEncoder(rw).Encode(&ret)
		return
	}

	defer func(con *pgx.Conn, ctx context.Context) {
		err := con.Close(ctx)
		if err != nil {
			log.Default().Printf("con.Close error : [%v]\n", err)
		}
	}(con, ctx)

	queries := sqlc.New(con)
	person, errCreateQuery := queries.CreatePerson(ctx, sqlc.CreatePersonParams{
		Personname:   *l.req.Body.PersonName,
		Email:        *l.req.Body.Email,
		PasswordHash: *l.req.Body.Password, // todo store password hash
	})
	if errCreateQuery != nil {
		log.Default().Printf("queries.CreatePerson error : [%v]\n", errCreateQuery)
		errStr := fmt.Sprintf("queries.CreatePerson error : [%v]", errCreateQuery)
		rw.WriteHeader(http.StatusInternalServerError)
		var ret models.Error
		ret.Code = http.StatusInternalServerError
		ret.Message = &errStr
		json.NewEncoder(rw).Encode(&ret)
		return
	}
	rw.WriteHeader(http.StatusOK)
	ret := models.RegistrationResponse{
		Email:      &person.Email,
		ID:         int64(person.ID),
		PersonName: &person.Personname,
	}
	json.NewEncoder(rw).Encode(&ret)
}

func RegisterHandlerFunc(params auth.RegisterParams) middleware.Responder {
	loginHandlerRet := &RegisterHandler{}
	loginHandlerRet.BindRequest(params)
	return loginHandlerRet
}
