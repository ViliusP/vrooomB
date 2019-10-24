package requests

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
	Requests        []Request `json:"requests,omitempty"`
}

type Status struct {
	StatusID int    `json:"status_ID"`
	Name     string `json:"name"`
}

type City struct {
	CityID int    `json:"cityID,omitempty"`
	Name   string `json:"name"`
}

type User struct {
	UserID   int    `json:"userID"`
	FName    string `json:"first_name"`
	LName    string `json:"last_name"`
	Phone    string `json:"phone_number"`
	PhotoURL string `json:"photo_URL"`
}

type Request struct {
	RequestID  int    `json:"requestID"`
	SubmitDate string `json:"submit_date"`
	Info       string `json:"info"`
	Requester  *User  `json:"request_user,omitempty"`
	Trip       *Trip  `json:"trip,omitempty"`
}
