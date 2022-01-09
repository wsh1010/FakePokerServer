package manager

import "net/http"

type User_Info struct {
	ID          string `json:"id"`
	Password    string `json:"password"`
	Describe    string `json:"describe"`
	Email       string `json:"email"`
	Phonenumber string `json:"phonenumber"`
	Name        string `json:"name"`
}

func Handler_userInfo() func(http.ResponseWriter, *http.Request) {
	return func(http.ResponseWriter, *http.Request) {

	}
}
