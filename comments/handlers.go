package comments

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"../util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

func GetComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tripID := mux.Vars(r)["id_TRIP"]
	personID := mux.Vars(r)["id_USER"]
	query := `
	SELECT trips.id_TRIP,trips.create_date,trips.departure_date,trips.cost_per_person,trips.space,trips.info,
	dep_city.id_CITY as id_CITY_dep, dep_city.name as name_dep,
	dest_city.id_CITY as id_CITY_dest, dest_city.name as name_dest,
	tripP.id_PERSON, tripP.first_name, tripP.last_name, tripP.phone_number, tripP.photo_URL,
    comments.id_COMMENT, comments.text,comments.submit_date,
	commentP.id_PERSON, commentP.first_name, commentP.last_name, commentP.phone_number, commentP.photo_URL
	FROM trips
	INNER JOIN CITIES as dep_city
	ON dep_city.id_CITY=fk_departure_CITY
	INNER JOIN CITIES as dest_city
	ON dest_city.id_CITY=fk_destination_CITY
	INNER JOIN people tripP 
	ON tripP.id_PERSON=fk_PERSONid_PERSON
    INNER JOIN comments 
    ON comments.fk_TRIP = trips.id_TRIP AND comments.fk_PERSON_TRIP = tripP.id_PERSON
    INNER JOIN people commentP 
    ON commentP.id_PERSON = comments.fk_PERSON_COMMENT
	WHERE tripP.isDeleted = 0 AND trips.id_TRIP = ? AND tripP.id_PERSON = ?
	`
	result, err := util.DB.Query(query, tripID, personID)
	if err != nil {
		fmt.Printf("Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var user User
	var trip Trip
	var departure City
	var destination City
	var comments []Comment
	for result.Next() {
		var userComment User
		var comment Comment

		err := result.Scan(
			&trip.TripID, &trip.CreationDate, &trip.DepartureDate, &trip.CostPerPerson, &trip.Space, &trip.Info,
			&departure.CityID, &departure.Name,
			&destination.CityID, &destination.Name,
			&user.UserID, &user.FName, &user.LName, &user.Phone, &user.PhotoURL,
			&comment.ID, &comment.Text, &comment.Date,
			&userComment.UserID, &userComment.FName, &userComment.LName, &userComment.Phone, &userComment.PhotoURL)
		if err != nil {
			fmt.Printf("Error: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		comment.Commentator = &userComment
		comments = append(comments, comment)

	}
	if trip.TripID == 0 && trip.CreationDate == "" {
		fmt.Printf("Error: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	trip.TripOwner = &user
	trip.DepartureCity = departure
	trip.DestinationCity = destination
	trip.TripComments = append(trip.TripComments, comments...)
	tripJSON, err := json.Marshal(trip)
	if err != nil {
		fmt.Printf("Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(tripJSON))

}

func GetCommentsByID(w http.ResponseWriter, r *http.Request) {
}

func DeleteCommentByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	commentatorID := parseID(r)
	tripID := mux.Vars(r)["id_TRIP"]
	personID := mux.Vars(r)["id_USER"]
	commentID := mux.Vars(r)["id_COMMENT"]
	query := `
	DELETE FROM comments WHERE id_COMMENT = ? AND fk_PERSON_COMMENT = ? AND fk_TRIP = ?	AND fk_PERSON_TRIP = ? `
	results, err := util.DB.Exec(query, commentID, commentatorID, tripID, personID)
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

func InsertComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	commentatorID := parseID(r)
	tripID := mux.Vars(r)["id_TRIP"]
	personID := mux.Vars(r)["id_USER"]

	query := `
	INSERT INTO comments(fk_TRIP, fk_PERSON_COMMENT, text, submit_date, fk_PERSON_TRIP) VALUES (?,?,?,?,?)
	`
	decoder := json.NewDecoder(r.Body)
	var comment Comment
	if err := decoder.Decode(&comment); err != nil || comment.Text == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	results, err := util.DB.Exec(query, tripID, commentatorID, comment.Text, time.Now().Format("2006-01-02 15:04:05"), personID)
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

func parseID(r *http.Request) string {
	user := r.Context().Value("user")
	return user.(*jwt.Token).Claims.(jwt.MapClaims)["id"].(string)
}
