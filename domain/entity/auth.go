package entity

import "github.com/dgrijalva/jwt-go"

// Authentication data model
type Authentication struct {
	Token      string
	ValidTime  int64
	ServerTime int64
}

type TokenData struct {
	jwt.StandardClaims
	AccountUUID string
}
