// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"github.com/go-openapi/loads"
	"github.com/yunusemreayhan/goAuthMicroService/internal/handler/auth_handlers"
	"github.com/yunusemreayhan/goAuthMicroService/internal/handler/person_handlers"
	"github.com/yunusemreayhan/goAuthMicroService/internal/security"
	"github.com/yunusemreayhan/goAuthMicroService/restapi/operations/auth"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/yunusemreayhan/goAuthMicroService/restapi/operations"
	"github.com/yunusemreayhan/goAuthMicroService/restapi/operations/person"
)

//go:generate swagger generate server --target ../../authMicroService2 --name AuthMicroService --spec ../swagger.yml --principal interface{}

func configureFlags(api *operations.AuthMicroServiceAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func getAPI() (*operations.AuthMicroServiceAPI, error) {
	swaggerSpec, err := loads.Analyzed(SwaggerJSON, "")
	if err != nil {
		return nil, err
	}
	api := operations.NewAuthMicroServiceAPI(swaggerSpec)
	return api, nil
}

func GetAPIHandler() (http.Handler, error) {
	api, err := getAPI()
	if err != nil {
		return nil, err
	}
	h := configureAPI(api)
	err = api.Validate()
	if err != nil {
		return nil, err
	}
	return h, nil
}

func configureAPI(api *operations.AuthMicroServiceAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	// Applies when the "X-API-Key" header is set
	api.APIKeyAuthAuth = security.ApiKeyAuthHandler
	// Applies when the Authorization header is set with the Basic scheme
	api.BasicAuthAuth = security.BasicAuthHandler

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()

	api.AuthRegisterHandler = auth.RegisterHandlerFunc(auth_handlers.RegisterHandlerFunc)
	api.AuthLoginHandler = auth.LoginHandlerFunc(auth_handlers.LoginHandlerFunc)
	api.AuthVerifyHandler = auth.VerifyHandlerFunc(auth_handlers.VerifyHandlerFunc)

	api.PersonDestroyOneHandler = person.DestroyOneHandlerFunc(person_handlers.DestroyOneHandlerFunc)
	api.PersonUpdateOneHandler = person.UpdateOneHandlerFunc(person_handlers.UpdateOneHandlerFunc)

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
