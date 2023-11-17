package auth_handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/yunusemreayhan/goAuthMicroService/models"
	"github.com/yunusemreayhan/goAuthMicroService/restapi/operations/auth"
)

type VerifyHandler struct {
	middleware.Responder
}

// curl -X GET "http://localhost/api/verify" -H "oauth_token:YOUR_TOKEN"

func (l *VerifyHandler) WriteResponse(rw http.ResponseWriter, pr runtime.Producer) {
	// if authentication is passed successfully, it will return "OK"
	rw.WriteHeader(http.StatusOK)
	ret := models.VerifyResponse{
		Nr200: "OK",
	}
	json.NewEncoder(rw).Encode(&ret)
}

func VerifyHandlerFunc(params auth.VerifyParams, principal interface{}) middleware.Responder {
	loginHandlerRet := &VerifyHandler{}
	return loginHandlerRet
}
