package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const keyPass = "leequ5tiek4zoo1mee2cheHa7ohmaigaiF9maigeig"

type appMetadata struct {
	Roles []string `json:"roles"`
}

type customClaims struct {
	jwt.StandardClaims
	AppMetadata appMetadata `json:"app_metadata"`
}

func main() {
	keyBytes, err := ioutil.ReadFile("keys/private.pem")
	if err != nil {
		panic(err)
	}

	pem, _ := pem.Decode(keyBytes)
	decoded, err := x509.DecryptPEMBlock(pem, []byte(keyPass))

	key, err := x509.ParsePKCS1PrivateKey(decoded)
	if err != nil {
		panic(err)
	}

	claims := &customClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			Issuer:    "test",
			Audience:  "demo",
			Subject:   "marcus",
		},
		AppMetadata: appMetadata{
			Roles: []string{"admins"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(key)

	fmt.Printf("Singed string: %s\n", signed)
}
