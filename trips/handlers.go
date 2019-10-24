package trips

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"../util"
)

func GetTrips(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var trips []Trip
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

	query := `
	SELECT trips.id_TRIP,trips.create_date,trips.departure_date,trips.cost_per_person,trips.space,trips.info,
	dep_city.id_CITY as id_CITY_dep, dep_city.name as name_dep,
	dest_city.id_CITY as id_CITY_dest, dest_city.name as name_dest,
	people.id_PERSON, people.first_name, people.last_name, people.phone_number, people.photo_URL
	FROM trips
	INNER JOIN CITIES as dep_city
	ON dep_city.id_CITY=fk_departure_CITY
	INNER JOIN CITIES as dest_city
	ON dest_city.id_CITY=fk_destination_CITY
	INNER JOIN people 
	ON people.id_PERSON=fk_PERSONid_PERSON
	WHERE people.isDeleted = 0 LIMIT ?,?
	`
	result, err := util.DB.Query(query, start, count)
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
		err := result.Scan(
			&trip.TripID, &trip.CreationDate, &trip.DepartureDate, &trip.CostPerPerson, &trip.Space, &trip.Info,
			&departure.CityID, &departure.Name,
			&destination.CityID, &destination.Name,
			&user.UserID, &user.FName, &user.LName, &user.Phone, &user.PhotoURL)
		if err != nil {
			fmt.Printf("Error: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		trip.TripOwner = &user
		trip.DepartureCity = departure
		trip.DestinationCity = destination
		trips = append(trips, trip)
	}
	tripsJSON, err := json.Marshal(trips)
	if err != nil {
		fmt.Printf("Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(tripsJSON))
}

func GetUserTrips(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	var trips []Trip

	query := `
	SELECT trips.id_TRIP,trips.create_date,trips.departure_date,trips.cost_per_person,trips.space,trips.info,
	dep_city.id_CITY as id_CITY_dep, dep_city.name as name_dep,
	dest_city.id_CITY as id_CITY_dest, dest_city.name as name_dest
	FROM trips
	INNER JOIN CITIES as dep_city
	ON dep_city.id_CITY=fk_departure_CITY
	INNER JOIN CITIES as dest_city
	ON dest_city.id_CITY=fk_destination_CITY
    INNER JOIN people
	ON people.id_PERSON=fk_PERSONid_PERSON
	WHERE people.isDeleted = 0 AND people.id_PERSON=?
	`

	result, err := util.DB.Query(query, id)
	if err != nil {
		fmt.Printf("Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for result.Next() {
		var trip Trip
		var departure City
		var destination City
		err := result.Scan(
			&trip.TripID, &trip.CreationDate, &trip.DepartureDate, &trip.CostPerPerson, &trip.Space, &trip.Info,
			&departure.CityID, &departure.Name,
			&destination.CityID, &destination.Name,
		)
		if err != nil {
			fmt.Printf("Error: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		trip.DepartureCity = departure
		trip.DestinationCity = destination
		trips = append(trips, trip)
	}
	tripsJSON, err := json.Marshal(trips)
	if err != nil {
		fmt.Printf("Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(tripsJSON))
}

func UpdateTripByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	decoder := json.NewDecoder(r.Body)
	var trip Trip
	if err := decoder.Decode(&trip); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	query := `UPDATE trips SET departure_date=?,cost_per_person=?,space=?,info=? WHERE trips.id_TRIP=?`
	if trip.DepartureDate == "" || trip.CostPerPerson < 0 || trip.Space < 0 || trip.Info == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	results, err := util.DB.Exec(query, trip.DepartureDate, trip.CostPerPerson, trip.Space, trip.Info, id)
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

func DeleteTripByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	query := `DELETE FROM trips WHERE trips.id_TRIP=?`
	results, err := util.DB.Exec(query, id)
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

func CreateTrip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	dateFormat := "2006-01-02 15:04:05" // YYYY:DD:MM HH:MM:SS
	query := `INSERT INTO trips(create_date, departure_date, cost_per_person, space, info, fk_departure_CITY, fk_PERSONid_PERSON, fk_destination_CITY) 
	VALUES (?,?,?,?,?,?,?,?)`
	decoder := json.NewDecoder(r.Body)
	var trip Trip
	if err := decoder.Decode(&trip); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	t1, e := time.Parse(dateFormat, trip.DepartureDate)
	fmt.Println(t1)
	fmt.Println(t1.Sub(time.Now()))
	if trip.Space < 0 || trip.CostPerPerson < 0 || t1.Sub(time.Now()) < 0 || e != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err := util.DB.Exec(query, time.Now().Format(dateFormat),
		trip.DepartureDate, trip.CostPerPerson, trip.Space, trip.Info,
		trip.DepartureCity.CityID, trip.TripOwner.UserID, trip.DestinationCity.CityID)

	if err != nil {
		fmt.Printf("Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
