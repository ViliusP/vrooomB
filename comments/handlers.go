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

func GetComment(w http.ResponseWriter, r *http.Request) {
}

func GetCommentsByID(w http.ResponseWriter, r *http.Request) {
}

func DeleteCommentByID(w http.ResponseWriter, r *http.Request) {

}

func InsertComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	commentatorID := parseID(r)
	tripID := mux.Vars(r)["id_TRIP"]
	query := `
	INSERT INTO comments(fk_TRIP, id_PERSON_COMMENT, text, submit_date) VALUES (?,?,?,?)
	`
	decoder := json.NewDecoder(r.Body)
	var comment Comment
	if err := decoder.Decode(&comment); err != nil || comment.Text == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	results, err := util.DB.Exec(query, tripID, commentatorID, comment.Text, time.Now().Format("2006-01-02 15:04:05"))
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
