package comments

type Comment struct {
	ID          string
	Date        string
	Text        string
	Commentator *User `json:"commentator,omitempty"`
}

type User struct {
	UserID   string `json:"userID"`
	FName    string `json:"first_name"`
	LName    string `json:"last_name"`
	Phone    string `json:"phone_number, omitempty"`
	PhotoURL string `json:"photo_URL"`
}

type Trip struct {
	TripID          int       `json:"tripID"`
	CreationDate    string    `json:"creation_date"`
	DepartureDate   string    `json:"departure_date"`
	CostPerPerson   int       `json:"cost"`
	Space           int       `json:"space"`
	Info            string    `json:"info"`
	DepartureCity   City      `json:"departure_city"`
	DestinationCity City      `json:"destination_city"`
	TripOwner       *User     `json:"user,omitempty"`
	TripComments    []Comment `json:"comments"`
}

type City struct {
	CityID string `json:"cityID,omitempty"`
	Name   string `json:"name"`
}
