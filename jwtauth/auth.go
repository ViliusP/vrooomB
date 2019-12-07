package jwtauth

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"../util"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

const (
	privKeyPath = "auth_keys/app.rsa"
	pubKeyPath  = "auth_keys/app.rsa.pub"
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

func init() {
	var err error

	SignKeyByte, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		log.Fatal("Error reading private key")
		return
	}
	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(SignKeyByte)
	if err != nil {
		log.Fatalf("[initKeys]: %s\n", err)
	}
	VerifyKeyByte, err := ioutil.ReadFile(pubKeyPath)
	if err != nil {
		log.Fatal("Error reading public key")
		return
	}
	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(VerifyKeyByte)
	if err != nil {
		log.Fatalf("[initKeys]: %s\n", err)
		panic(err)
	}
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var creds Credentials
	var user User
	exists := false
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	q := `SELECT id_PERSON, email, password, isDeleted FROM people WHERE email = ? AND password=?`
	result, err := util.DB.Query(q, creds.Email, creds.Password)
	if err != nil {
		fmt.Printf("Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Get the expected password from our in memory map

	for result.Next() {
		exists = true
		err := result.Scan(&user.id, &user.email, &user.password, &user.isDeleted)
		if err != nil {
			fmt.Printf("Error: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if user.password != creds.Password || user.isDeleted == 1 || !exists {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	signer := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), jwt.MapClaims{
		"id":    user.id,
		"email": user.email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
		"iss":   "localhostas",
	})

	tokenString, err := signer.SignedString(signKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error signing token: %v\n", err)
	}

	//create a token instance using the token string
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", "Bearer "+tokenString)
	w.WriteHeader(http.StatusNoContent)

}

func CheckJWT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func AuthMiddleware(next http.Handler) http.Handler {
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return verifyKey, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})
	return jwtMiddleware.Handler(next)
}
