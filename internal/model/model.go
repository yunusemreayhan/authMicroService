package model

type JWTIdentity struct {
	Id       int32  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type RegistrationRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegistrationResponse struct {
	Id       int32  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type VerifyVoucherRequest struct {
	Token string `json:"voucher"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
