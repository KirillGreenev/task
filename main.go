package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type User struct {
	Id        int     `json:"id"`
	Email     string  `json:"email,omitempty"`
	Amount    int     `json:"amount"`
	Profile   Profile `json:"profile,omitempty"`
	Password  string  `json:"-"`
	Username  string  `json:"username"`
	CreatedAt string  `json:"createdAt"`
	CreatedBy string  `json:"createdBy"`
}

type Profile struct {
	Avatar     string `json:"avatar,omitempty"`
	LastName   string `json:"lastName,omitempty"`
	FirstName  string `json:"firstName,omitempty"`
	StaticData string `json:"-"`
}

func main() {
	http.HandleFunc("/user", GetUser)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func GetUser(w http.ResponseWriter, _ *http.Request) {
	users, err := GetUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, user := range users {
		if user.Amount > 50000 {
			user.Email = ""
			user.Profile = Profile{}
		}
	}

	json.NewEncoder(w).Encode(users)
}

func GetUsers() ([]*User, error) {
	r, err := http.NewRequest("GET", "http://83.136.232.77:8091/users", nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var users []*User
	if err = json.NewDecoder(res.Body).Decode(&users); err != nil {
		return nil, err
	}

	return users, nil
}
