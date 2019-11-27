package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"../util"
	"github.com/gorilla/mux"
)

//GetUsers get all users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var users []User
	count, err1 := strconv.Atoi(r.FormValue("limit"))
	start, err2 := strconv.Atoi(r.FormValue("offset"))
	if count == 0 && start == 0 {
		count = 10
		start = 0
	}
	if count <= 0 || start < 0 || err1 != nil || err2 != nil {
		count = 10
		start = 0
		// w.WriteHeader(http.StatusBadRequest)
		// return
	}

	result, err := util.DB.Query("SELECT `id_PERSON`,`first_name`,`last_name`,`phone_number`,`email`,`registration_date`,`photo_URL` FROM `people` WHERE people.isDeleted = 0 LIMIT ?,?", start, count)
	if err != nil {
		fmt.Printf("Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for result.Next() {
		var user User
		err := result.Scan(&user.UserID, &user.FName, &user.LName, &user.Phone, &user.Email, &user.RegistrationDate, &user.PhotoURL)
		if err != nil {
			fmt.Printf("Error: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}
	usersJSON, err := json.Marshal(users)
	if err != nil {
		fmt.Printf("Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(usersJSON))

}

//GetUserByID ...
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if id < 0 || err != nil {
		fmt.Printf("Error: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result, err := util.DB.Query("SELECT `id_PERSON`,`first_name`,`last_name`,`phone_number`,`email`,`registration_date`,`photo_URL` FROM `people` WHERE people.id_PERSON = ? AND people.isDeleted = 0", id)
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
	userJSON, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	if user.UserID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, string(userJSON))

}

//UpdateUserByID ...
func UpdateUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if id < 0 || err != nil {
		fmt.Printf("Error: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var user User
	if err := decoder.Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if user.FName == "" || user.LName == "" || user.Phone == "" || user.PhotoURL == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	results, err := util.DB.Exec("UPDATE `people` SET `first_name`=?,`last_name`=?,`phone_number`=?,`photo_URL`=? WHERE `id_PERSON`=?", user.FName, user.LName, user.Phone, user.PhotoURL, params["id"])
	if err != nil {
		fmt.Printf("Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	RowsAffected, _ := results.RowsAffected()
	if RowsAffected <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

//DeleteUserByID should I check if there is user in DB before deleting???
func DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if id < 0 || err != nil {
		fmt.Printf("Error: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	results, err := util.DB.Exec("UPDATE `people` SET `isDeleted`=1 WHERE people.id_PERSON=? AND people.isDeleted = 0", params["id"])
	if err != nil {
		fmt.Printf("Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	RowsAffected, _ := results.RowsAffected()
	if RowsAffected <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
}
