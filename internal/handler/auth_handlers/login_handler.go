package auth_handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/yunusemreayhan/goAuthMicroService/db/sqlc"
	"github.com/yunusemreayhan/goAuthMicroService/internal/config"
	"github.com/yunusemreayhan/goAuthMicroService/internal/key"
	"github.com/yunusemreayhan/goAuthMicroService/internal/model"
	"github.com/yunusemreayhan/goAuthMicroService/models"
	"github.com/yunusemreayhan/goAuthMicroService/restapi/operations/auth"
	"net/http"
)

type LoginHandler struct {
	middleware.Responder
	req    auth.LoginParams
	person *sqlc.Person
}

func (l *LoginHandler) BindRequest(req interface{}) {
	l.person = req.(*sqlc.Person)
}

func (l *LoginHandler) WriteResponse(rw http.ResponseWriter, pr runtime.Producer) {
	var loginResponse models.LoginResponse

	var errCreatingToken error
	loadedKey, mError := key.LoadPrivateKey(config.GetConfig().DefaultKeyPath)
	if mError != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		errStr := fmt.Sprintf("key.LoadPrivateKey mError : [%v]", mError)
		ret := models.Error{
			Code:    0,
			Message: &errStr,
		}
		json.NewEncoder(rw).Encode(&ret)
		return
	}
	if l.person == nil {
		rw.WriteHeader(http.StatusBadRequest)
		errStr := fmt.Sprintf("can not load person from request")
		ret := models.Error{
			Code:    0,
			Message: &errStr,
		}
		json.NewEncoder(rw).Encode(&ret)
		return
	}
	token, errCreatingToken := key.CreateJwtTokenWithIdentity(model.JWTIdentity{
		Id:       l.person.ID,
		Username: l.person.Personname,
		Email:    l.person.Email,
	}, loadedKey)
	if errCreatingToken != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		errStr := fmt.Sprintf("key.CreateJwtTokenWithIdentity mError : [%v]", errCreatingToken)
		ret := models.Error{
			Code:    0,
			Message: &errStr,
		}
		json.NewEncoder(rw).Encode(&ret)
		return
	}
	loginResponse.Token = &token

	json.NewEncoder(rw).Encode(&loginResponse)
}

func LoginHandlerFunc(params auth.LoginParams, principal interface{}) middleware.Responder {
	loginHandlerRet := &LoginHandler{}
	loginHandlerRet.BindRequest(principal)
	return loginHandlerRet
}
