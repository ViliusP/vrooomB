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
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	q := `SELECT id_PERSON, email, password, isDeleted FROM people WHERE email = ?`
	result, err := util.DB.Query(q, creds.Email)
	if err != nil {
		fmt.Printf("Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Get the expected password from our in memory map
	var user User
	for result.Next() {
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
	if user.password != creds.Password || user.isDeleted == 1 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	signer := jwt.New(jwt.GetSigningMethod("RS256"))
	claims := make(jwt.MapClaims)
	claims["iss"] = "admin"
	claims["exp"] = time.Now().Add(time.Minute * 20).Unix()
	claims["UserInfo"] = struct {
		email string
		id    string
	}{user.email, user.id}
	signer.Claims = claims
	tokenString, err := signer.SignedString(signKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error signing token: %v\n", err)
	}

	//create a token instance using the token string
	response := Token{tokenString}
	JSONResponse(response, w)
}

func JSONResponse(response interface{}, w http.ResponseWriter) {

	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
