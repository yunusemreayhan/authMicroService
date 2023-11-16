package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
	"github.com/yunusemreayhan/goAuthMicroService/internal/config"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/jackc/pgx/v5"
	dbutils "github.com/yunusemreayhan/goAuthMicroService/db"
	"github.com/yunusemreayhan/goAuthMicroService/db/sqlc"
	"github.com/yunusemreayhan/goAuthMicroService/internal/key"
	"github.com/yunusemreayhan/goAuthMicroService/internal/model"
)

// Authentication Microservice

// @Summary Person Registration
// @Description Register a new person
// @ID register-person
// @Accept json
// @Produce json
// @Param person body RegistrationRequest true "Person registration data"
// @Success 201 {object} RegistrationResponse
// @Router /api/register [post]
func RegisterPerson(c *fiber.Ctx) error {
	// Handle person registration
	// try parsing request as json
	var request model.RegistrationRequest
	var username, email, password string
	res := json.Unmarshal(c.Request().Body(), &request)
	defer func(request *fasthttp.Request) {
		err := request.CloseBodyStream()
		if err != nil {
			log.Default().Printf("request.CloseBodyStream error : [%v]\n", err)
		}
	}(c.Request())

	if res != nil {
		username = c.FormValue("username")
		password = c.FormValue("password")
		email = c.FormValue("email")
	} else {
		username = request.Username
		password = request.Password
		email = request.Email
	}

	if username == "" || password == "" || email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username, password and email are required!",
		})
	}

	ctx := context.Background()

	con, err := pgx.Connect(ctx, dbutils.GetSQLDSN())
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "DB not reachable!",
		})
	}

	defer func(con *pgx.Conn, ctx context.Context) {
		err := con.Close(ctx)
		if err != nil {
			log.Default().Printf("con.Close error : [%v]\n", err)
		}
	}(con, ctx)

	queries := sqlc.New(con)
	person, errDB := queries.CreatePerson(ctx, sqlc.CreatePersonParams{
		Personname:   username,
		Email:        email,
		PasswordHash: password,
	})

	if errDB != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": fmt.Sprintf("DB person creation error [%s]", errDB.Error()),
		})
	}

	return c.JSON(model.RegistrationResponse{
		Id:       person.ID,
		Username: username,
		Email:    email,
	})
}

// LoginPerson
// @Summary Person Login
// @Description Log in a person and receive a JWT token
// @ID login-person
// @Accept json
// @Produce json
// @Param person body LoginRequest true "Person login data"
// @Success 200 {string} string "JWT Token"
// @Router /api/login [post]
func LoginPerson(c *fiber.Ctx) error {
	// try parsing request as json
	var request model.LoginRequest
	var username, email, password string
	res := json.Unmarshal(c.Request().Body(), &request)
	defer func(request *fasthttp.Request) {
		err := request.CloseBodyStream()
		if err != nil {
			log.Default().Printf("request.CloseBodyStream error : [%v]\n", err)
		}
	}(c.Request())

	if res != nil {
		// Handle person login
		username = c.FormValue("username")
		email = c.FormValue("email")
		password = c.FormValue("password")
	} else {
		username = request.Username
		email = request.Email
		password = request.Password
	}

	if username == "" && email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing username and email, please provide one!",
		})
	}

	if password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing password, please provide one!",
		})
	}

	ctx := context.Background()

	con, err := pgx.Connect(ctx, dbutils.GetSQLDSN())
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "DB not reachable!",
		})
	}

	defer func(con *pgx.Conn, ctx context.Context) {
		err := con.Close(ctx)
		if err != nil {
			log.Default().Printf("con.Close error : [%v]\n", err)
		}
	}(con, ctx)

	queries := sqlc.New(con)

	var voucherOwnerPerson sqlc.Person
	var errDB error
	if email != "" {
		voucherOwnerPerson, errDB = queries.GetPersonByEmail(ctx, email)
		if errDB != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": fmt.Sprintf("DB person with email [%s] not exist error [%s]", email, errDB.Error()),
			})
		}
	}

	if username != "" {
		voucherOwnerPerson, errDB = queries.GetPersonByPersonName(ctx, username)
		if errDB != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": fmt.Sprintf("DB person with username [%s] not exist error [%s]", username, errDB.Error()),
			})
		}
	}

	if voucherOwnerPerson.PasswordHash != password {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"identifier": model.JWTIdentity{Id: voucherOwnerPerson.ID, Username: voucherOwnerPerson.Personname, Email: voucherOwnerPerson.Email},
		"admin":      true,
		"exp":        time.Now().Add(time.Duration(config.GetConfig().TokenTimeoutSeconds) * time.Second).Unix(),
	}

	// Create voucher
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Generate encoded voucher and send it as response.
	privateKey, err := key.LoadPrivateKey(key.DefaultKeyPath)
	if err != nil {
		log.Printf("key.LoadPrivateKey: [%v]\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	t, err := token.SignedString(privateKey)
	if err != nil {
		log.Printf("token.SignedString: [%v]\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"voucher": t})
}

// VerifyToken
// @Summary Token Verification
// @Description Verify the provided JWT token
// @ID verify-token
// @Produce json
// @Param Authorization header string true "JWT Token"
// @Success 200 {string} string "Access granted"
// @Failure 401 {string} string "Access denied"
// @Router /api/verify [get]
func VerifyToken(c *fiber.Ctx) error {
	var request model.VerifyVoucherRequest
	res := json.Unmarshal(c.Request().Body(), &request)
	defer func(request *fasthttp.Request) {
		err := request.CloseBodyStream()
		if err != nil {
			log.Default().Printf("request.CloseBodyStream error : [%v]\n", err)
		}
	}(c.Request())

	if res != nil {
		log.Printf("json.Unmarshal error : [%v] failed to parse request [%v]\n", res, string(c.Request().Body()))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf(`Was not able to parse request json json.Unmarshal: [%s]`, res.Error())})
	}

	// Generate encoded voucher and send it as response.
	privateKey, err := key.LoadPrivateKey(key.DefaultKeyPath)
	if err != nil {
		log.Printf("key.LoadPrivateKey: [%v]\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not load private key"})
	}

	token, err := jwt.Parse(request.Token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: [%v]", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return privateKey.Public(), nil
	})

	if err == nil {
		if token == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}
		return c.SendStatus(fiber.StatusOK)
	}

	switch {
	case errors.Is(err, jwt.ErrTokenMalformed):
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "That's not even a token"})
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid signature"})
	case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Timing is everything"})
	default:
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": fmt.Sprintf("Couldn't handle this token: [%s]", err)})
	}
}

// Person Database Microservice

// UpdatePerson
// @Summary Update Person Information
// @Description Update person information in the database
// @ID create-person
// @Accept json
// @Produce json
// @Param person body PersonUpdateRequest true "Person update data"
// @Success 201 {string} string "Person information stored"
// @Router /api/person [post]
func UpdatePerson(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["identifier"].(map[string]interface{})["username"].(string)

	ctx := context.Background()

	con, err := pgx.Connect(ctx, dbutils.GetSQLDSN())
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "DB not reachable!",
		})
	}

	defer func(con *pgx.Conn, ctx context.Context) {
		err := con.Close(ctx)
		if err != nil {
			log.Default().Printf("con.Close error : [%v]\n", err)
		}
	}(con, ctx)

	queries := sqlc.New(con)

	voucherOwnerPerson, errQueryUser := queries.GetPersonByPersonName(ctx, username)
	if errQueryUser != nil {
		log.Default().Printf("DB person with username [%s]  not exist error [%s]\n", username, errQueryUser.Error())
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": fmt.Sprintf("DB person with username [%s]  not exist error [%s]", username, errQueryUser.Error()),
		})
	}

	if c.Method() == "PUT" {
		// Handle person login
		password := c.FormValue("password")

		if password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Missing password, please provide one!",
			})
		}

		_, errDB := queries.UpdatePersonPasswordHashById(ctx, sqlc.UpdatePersonPasswordHashByIdParams{
			ID:           voucherOwnerPerson.ID,
			PasswordHash: password,
		})

		if errDB != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to update person password [%s]", errDB.Error()),
			})
		}

		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"error": "",
		})
	}

	return c.Status(fiber.StatusMethodNotAllowed).JSON(fiber.Map{
		"error": "Method not allowed!",
	})
}
