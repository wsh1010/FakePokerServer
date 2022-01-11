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
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			//회원 저장.
		case http.MethodPut:
			//회원 수정.
		case http.MethodDelete:
			//회원 탈퇴.
		case http.MethodGet:
			//회원 조회.
		default:
			http.Error(w, "Invalid request method.", http.StatusBadRequest)
		}
	}
}

func Handler_login() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request){
		switch r.Method {
		case http.MethodGet:
						
		}
	}
}