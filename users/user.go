package users

type User struct {
	UserID string `json:"userID"`
	Email  string `json:"email"`
	//Password string `json:"password"`
	FName            string `json:"first_name"`
	LName            string `json:"last_name"`
	Phone            string `json:"phone_number"`
	PhotoURL         string `json:"photo_URL"`
	RegistrationDate string `json:"registration_date"`
	isDeleted        int
}
