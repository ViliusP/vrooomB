package requests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"../util"
	"github.com/gorilla/mux"
)

func GetUserRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var requests []Request
	id := mux.Vars(r)["id"]
	intID, err := strconv.Atoi(id)
	if err != nil || intID < 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	query := `
	SELECT 
	requests.id_REQUEST, requests.submit_date, requests.info,
	request_statuses.id_REQUEST_STATUS, request_statuses.name,
	trips.id_TRIP,trips.create_date,trips.departure_date,trips.cost_per_person,trips.space,trips.info,
	dep_city.id_CITY as id_CITY_dep, dep_city.name as name_dep,
	dest_city.id_CITY as id_CITY_dest, dest_city.name as name_des,
	people.id_PERSON, people.first_name, people.last_name, people.phone_number, people.photo_URL, people.isDeleted
	FROM requests 
	INNER JOIN request_statuses
	ON request_statuses.id_REQUEST_STATUS = requests.request_status
	LEFT JOIN trips
	ON trips.id_TRIP = requests.fk_TRIP
	INNER JOIN CITIES as dep_city
	ON dep_city.id_CITY=fk_departure_CITY
	INNER JOIN CITIES as dest_city
	ON dest_city.id_CITY=fk_destination_CITY
	INNER JOIN people
	ON people.id_PERSON = trips.fk_PERSONid_PERSON
	WHERE requests.fk_PERSON=? LIMIT ?,?
	`
	count, _ := strconv.Atoi(r.FormValue("limit"))
	start, _ := strconv.Atoi(r.FormValue("offset"))

	if count == 0 && start == 0 {
		count = 10
		start = 0
	}
	if count <= 0 || start < 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := util.DB.Query(query, id, start, count)
	if err != nil {
		fmt.Printf("Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for result.Next() {
		var user User
		var trip Trip
		var departure City
		var destination City
		var request Request
		var requestStatus Status
		isDeleted := 0
		err := result.Scan(
			&request.RequestID, &request.SubmitDate, &request.Info,
			&requestStatus.StatusID, &requestStatus.Name,
			&trip.TripID, &trip.CreationDate, &trip.DepartureDate, &trip.CostPerPerson, &trip.Space, &trip.Info,
			&departure.CityID, &departure.Name,
			&destination.CityID, &destination.Name,
			&user.UserID, &user.FName, &user.LName, &user.Phone, &user.PhotoURL, &isDeleted) //Last parameter is 'isDeleted'
		if err != nil {
			fmt.Printf("Error: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		trip.TripOwner = &user
		trip.DepartureCity = departure
		trip.DestinationCity = destination
		request.Trip = &trip
		requests = append(requests, request)
	}

	requestsJSON, err := json.Marshal(requests)
	if err != nil {
		fmt.Printf("Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(requestsJSON))
}

func GetTripRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	intID, err := strconv.Atoi(id)
	if err != nil || intID < 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	query := `
	SELECT 
	requests.id_REQUEST, requests.submit_date, requests.info,
	request_statuses.id_REQUEST_STATUS, request_statuses.name,
	trips.id_TRIP,trips.create_date,trips.departure_date,trips.cost_per_person,trips.space,trips.info,
	dep_city.id_CITY as id_CITY_dep, dep_city.name as name_dep,
	dest_city.id_CITY as id_CITY_dest, dest_city.name as name_des,
	people.id_PERSON, people.first_name, people.last_name, people.phone_number, people.photo_URL, people.isDeleted
	FROM requests 
	INNER JOIN request_statuses
	ON request_statuses.id_REQUEST_STATUS = requests.request_status
	INNER JOIN trips
	ON trips.id_TRIP = requests.fk_TRIP
	INNER JOIN CITIES as dep_city
	ON dep_city.id_CITY=fk_departure_CITY
	INNER JOIN CITIES as dest_city
	ON dest_city.id_CITY=fk_destination_CITY
	INNER JOIN people
	ON people.id_PERSON = requests.fk_PERSON
    WHERE requests.fk_TRIP=? LIMIT ?,?
	`
	count, _ := strconv.Atoi(r.FormValue("limit"))
	start, _ := strconv.Atoi(r.FormValue("offset"))

	if count == 0 && start == 0 {
		count = 10
		start = 0
	}
	if count <= 0 || start < 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := util.DB.Query(query, id, start, count)
	if err != nil {
		fmt.Printf("Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var requests []Request
	var trip Trip
	for result.Next() {
		var user User
		var departure City
		var destination City
		var request Request
		var requestStatus Status
		isDeleted := 0
		err := result.Scan(
			&request.RequestID, &request.SubmitDate, &request.Info,
			&requestStatus.StatusID, &requestStatus.Name,
			&trip.TripID, &trip.CreationDate, &trip.DepartureDate, &trip.CostPerPerson, &trip.Space, &trip.Info,
			&departure.CityID, &departure.Name,
			&destination.CityID, &destination.Name,
			&user.UserID, &user.FName, &user.LName, &user.Phone, &user.PhotoURL, &isDeleted) //Last parameter is 'isDeleted'
		if err != nil {
			fmt.Printf("Error: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		trip.DepartureCity = departure
		trip.DestinationCity = destination
		request.Requester = &user
		requests = append(requests, request)
	}
	trip.Requests = requests
	requestsJSON, err := json.Marshal(trip)
	if err != nil {
		fmt.Printf("Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(requestsJSON))
}

func DeleteRequestByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

}

func UpdateRequestByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

}

func CreateRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

}

func GetStatuses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var statuses []Status
	query := `
	SELECT id_REQUEST_STATUS, name
	FROM request_statuses
	`
	result, err := util.DB.Query(query)
	if err != nil {
		fmt.Printf("Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for result.Next() {
		var status Status
		err := result.Scan(&status.StatusID, &status.Name)
		if err != nil {
			fmt.Printf("Error: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		statuses = append(statuses, status)
	}
	statusesJSON, err := json.Marshal(statuses)
	if err != nil {
		fmt.Printf("Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(statusesJSON))
}
