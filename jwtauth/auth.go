package jwtauth

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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
	// claims := CustomClaims{
	// 	user.id,
	// 	user.email,
	// 	jwt.StandardClaims{
	// 		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	// 		Issuer:    "localhost",
	// 	},
	// }
	// signer := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), claims)
	signer := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), jwt.MapClaims{
		"id":    user.id,
		"email": user.email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
		"iss":   "localhostas",
	})

	// claims := make(jwt.MapClaims)
	// claims["iss"] = "localhost"
	// claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	// claims["UserInfo"] = UserInfo{user.id, user.email}
	// signer.Claims = claims
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

func AuthMiddleware(next http.Handler) http.Handler {
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return verifyKey, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})
	return jwtMiddleware.Handler(next)
}

// func AuthMiddleware(next http.Handler) http.Handler {
// 	return Handler(next)
// }

func Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Let secure process the request. If it returns an error,
		// that indicates the request should not continue.
		err := CheckJWT(w, r)

		// If there was an error, do not continue.
		if err != nil {
			return
		}
		h.ServeHTTP(w, r)
	})
}

// func RefreshToken(h http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		//	now := time.Now().Unix()
// 		extractor := func(token *jwt.Token) (interface{}, error) {
// 			return verifyKey, nil
// 		}
// 		tokenString := r.Context().Value("User").(*jwt.Token).Raw
// 		parsedToken, _ := jwt.Parse(tokenString, extractor)
// 		_ = parsedToken
// 		if true {
// 			userInfo := r.Context().Value("User").(*jwt.Token).Claims.(jwt.MapClaims)["UserInfo"]
// 			_ = userInfo
// 			//	token := newToken(userInfo)
// 			w.Header().Set("Autorization", "Bearer"+"token")
// 		}
// 		h.ServeHTTP(w, r)
// 	})
// }

func CheckJWT(w http.ResponseWriter, r *http.Request) error {
	extractor := func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	}

	if r.Method == "OPTIONS" {
		return nil
	}

	// Use the specified token extractor to extract a token from the request
	token, err := FromAuthHeader(r)

	// If the token is empty...
	if token == "" {
		// If we get here, the required token is missing
		errorMsg := "Required authorization token not found"
		OnError(w, r, errorMsg)
		return fmt.Errorf(errorMsg)
	}

	// Now parse the token
	parsedToken, err := jwt.Parse(token, extractor)
	formattedError := fmt.Sprintf("%v", err)
	// Check if there was an error in parsing...
	// fmt.Sprintf("%v", err) != "Token is expired" {
	if err != nil {
		OnError(w, r, err.Error())
		return fmt.Errorf("Error parsing token: %v", err)
	}

	if jwt.SigningMethodRS256 != nil && jwt.SigningMethodRS256.Alg() != parsedToken.Header["alg"] {
		message := fmt.Sprintf("Expected %s signing method but token specified %s",
			jwt.SigningMethodRS256.Alg(),
			parsedToken.Header["alg"])
		OnError(w, r, errors.New(message).Error())
		return fmt.Errorf("Error validating token algorithm: %s", message)
	}

	// Check if the parsed token is valid...
	if !parsedToken.Valid && formattedError != "Token is expired" {
		OnError(w, r, "The token isn't valid")
		return errors.New("Token is invalid")
	}

	// if formattedError == "Token is expired" {
	// 	newRequest := r.WithContext(context.WithValue(r.Context(), "isExpired", true))
	// 	*r = *newRequest
	// }
	// if formattedError != "Token is expired" {
	// 	newRequest := r.WithContext(context.WithValue(r.Context(), "isExpired", false))
	// 	*r = *newRequest
	// }

	// If we get here, everything worked and we can set the
	// user property in context.
	newRequest := r.WithContext(context.WithValue(r.Context(), "User", parsedToken))
	// Update the current request with the new context information.
	*r = *newRequest
	return nil
}

func FromAuthHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil // No error, just no token
	}

	// TODO: Make this a bit more robust, parsing-wise
	authHeaderParts := strings.Fields(authHeader)
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", errors.New("Authorization header format must be Bearer {token}")
	}

	return authHeaderParts[1], nil
}

func OnError(w http.ResponseWriter, r *http.Request, err string) {
	http.Error(w, err, http.StatusUnauthorized)
}
