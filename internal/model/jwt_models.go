package model

type JWTIdentity struct {
	Id       int32  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
