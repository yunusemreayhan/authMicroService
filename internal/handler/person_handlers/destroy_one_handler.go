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

type DestroyOneHandler struct {
	middleware.Responder
	person *sqlc.Person
}

func (l *DestroyOneHandler) BindRequest(req interface{}) {
	l.person = req.(*sqlc.Person)
}

func (l *DestroyOneHandler) WriteResponse(rw http.ResponseWriter, pr runtime.Producer) {
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
	errDeleteQuery := queries.DeletePersonByID(ctx, l.person.ID)
	if errDeleteQuery != nil {
		log.Default().Printf("queries.UpdatePerson error : [%v]\n", errDeleteQuery)
		errStr := fmt.Sprintf("queries.UpdatePerson error : [%v]", errDeleteQuery)
		rw.WriteHeader(http.StatusInternalServerError)
		var ret models.Error
		ret.Code = http.StatusInternalServerError
		ret.Message = &errStr
		json.NewEncoder(rw).Encode(&ret)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}

func DestroyOneHandlerFunc(params person.DestroyOneParams, principal interface{}) middleware.Responder {
	loginHandlerRet := &DestroyOneHandler{}
	loginHandlerRet.BindRequest(principal)
	return loginHandlerRet
}
