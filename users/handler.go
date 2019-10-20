package users

import (
	"encoding/json"
	"fmt"
	"net/http"

	"../util"
	"github.com/gorilla/mux"
)

var mockUser1 = &User{
	UserID:           "1",
	Email:            "mock1@mail.com",
	FName:            "FirstMock1",
	LName:            "LastMock1",
	Phone:            "+37066666666",
	PhotoURL:         "localhost/img/profiles/1",
	RegistrationDate: "1571499160",
	isDeleted:        0,
}

var mockUser2 = &User{
	UserID:           "2",
	Email:            "mock2@mail.com",
	FName:            "FirstMock2",
	LName:            "LastMock2",
	Phone:            "+3706666756",
	PhotoURL:         "localhost/img/profiles/2",
	RegistrationDate: "1571499223",
	isDeleted:        0,
}

var mockUser3 = &User{
	UserID:           "3",
	Email:            "mock3@mail.com",
	FName:            "FirstMock3",
	LName:            "LastMock3",
	Phone:            "+3706666655555",
	PhotoURL:         "localhost/img/profiles/3",
	RegistrationDate: "1571499217",
	isDeleted:        0,
}

var users = []User{*mockUser1, *mockUser2, *mockUser3}

//GetUsers get all users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	usersJSON, err := json.Marshal(users)
	if err != nil {
		fmt.Printf("Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Fprint(w, string(usersJSON))
}

//GetUserByID ...
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	result, err := util.DB.Query("SELECT `id_PERSON`,`first_name`,`last_name`,`phone_number`,`email`,`registration_date`,`photo_URL` FROM `people` WHERE people.id_PERSON = ? AND people.isDeleted = 0", params["id"])
	if err != nil {
		fmt.Printf("Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var user User
	for result.Next() {
		err := result.Scan(&user.UserID, &user.FName, &user.LName, &user.Phone, &user.Email, &user.RegistrationDate, &user.PhotoURL)
		if err != nil {
			fmt.Printf("Error: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	if user.UserID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		fmt.Fprintf(w, string(userJSON))
	}

}

func UpdateUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]
	userToSend := &User{}
	isFound := false
	for _, user := range users {
		if user.UserID == id && user.isDeleted == 0 {
			userToSend = &user
			isFound = true
			break
		}
	}
	userJSON, err := json.Marshal(userToSend)
	if err != nil {
		fmt.Printf("Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	} else if isFound {
		fmt.Fprintf(w, string(userJSON))
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

//DeleteUserByID should I check if there is user in DB before deleting???
func DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	_, err := util.DB.Query("UPDATE `people` SET `isDeleted`=1 WHERE people.id_PERSON=? AND people.isDeleted = 0", params["id"])
	if err != nil {
		fmt.Printf("Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func createUser(w http.ResponseWriter, r *http.Request) {

}
