package person_handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/jackc/pgx/v5"
	dbutils "github.com/yunusemreayhan/goAuthMicroService/db"
	"github.com/yunusemreayhan/goAuthMicroService/db/sqlc"
	"github.com/yunusemreayhan/goAuthMicroService/models"
	"github.com/yunusemreayhan/goAuthMicroService/restapi/operations/person"
	"log"
	"net/http"
)

type UpdateOneHandler struct {
	middleware.Responder
	updateOneRequest *models.Person
	person           *sqlc.Person
}

func (l *UpdateOneHandler) BindRequest(params person.UpdateOneParams, req interface{}) {
	l.person = req.(*sqlc.Person)
	l.updateOneRequest = params.Body
}

func (l *UpdateOneHandler) WriteResponse(rw http.ResponseWriter, pr runtime.Producer) {
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
	updatedPerson, errUpdateQuery := queries.UpdatePerson(ctx, sqlc.UpdatePersonParams{
		ID:           l.person.ID,
		Personname:   *l.updateOneRequest.PersonName,
		Email:        *l.updateOneRequest.Email,
		PasswordHash: *l.updateOneRequest.Password, // todo store hash
	})
	if errUpdateQuery != nil {
		log.Default().Printf("queries.UpdatePerson error : [%v]\n", errUpdateQuery)
		errStr := fmt.Sprintf("queries.UpdatePerson error : [%v]", errUpdateQuery)
		rw.WriteHeader(http.StatusInternalServerError)
		var ret models.Error
		ret.Code = http.StatusInternalServerError
		ret.Message = &errStr
		json.NewEncoder(rw).Encode(&ret)
		return
	}
	rw.WriteHeader(http.StatusOK)
	var putPersonResponse models.PersonWithoutPassword
	putPersonResponse.ID = int64(updatedPerson.ID)
	putPersonResponse.PersonName = &updatedPerson.Personname
	putPersonResponse.Email = &updatedPerson.Email

	json.NewEncoder(rw).Encode(&putPersonResponse)
}

func UpdateOneHandlerFunc(params person.UpdateOneParams, principal interface{}) middleware.Responder {
	loginHandlerRet := &UpdateOneHandler{}
	loginHandlerRet.BindRequest(params, principal)
	return loginHandlerRet
}
