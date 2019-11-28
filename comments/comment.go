package comments

type Comment struct {
	Text        string
	commentator *User `json:"commentator,omitempty"`
}

type User struct {
	UserID   string `json:"userID"`
	FName    string `json:"first_name"`
	LName    string `json:"last_name"`
	Phone    string `json:"phone_number"`
	PhotoURL string `json:"photo_URL"`
}
