package manager

import (
	"database/sql"
	"encoding/json"
	"fakepokerserver/module/db"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type User_Info struct {
	ID   string `json:"id"`
	Nick string `json:"nick"`
}

type User_CoinData struct {
	Coins int `json:"coins"`
}

type User_Status struct {
	Nick   string `json:"nick"`
	Status string `json:"status"`
	Room   string `json:"room_id"`
	Coins  int    `json:"coins"`
	Error  string `json:"error"`
}

func Handler_userInfo() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			//회원 저장.
			var body User_Info
			userID, exist := r.Header["Userid"]
			if !exist {
				log.Println("not userID")
				w.WriteHeader(http.StatusNotImplemented)
				return
			}
			respBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			json.Unmarshal(respBody, &body)
			saveUser(userID[0], body.Nick)
			w.WriteHeader(http.StatusOK)
		case http.MethodPut:
			//회원 수정.
		case http.MethodDelete:
			//회원 탈퇴.
		case http.MethodGet:
			//회원 조회.
			var responsebody User_Status
			userID, exist := r.Header["Userid"]
			if !exist {
				log.Println("not userID")
				responsebody.Error = "not userID"
				responsebodyBytes, _ := json.Marshal(responsebody)

				w.WriteHeader(http.StatusNotFound)
				w.Write(responsebodyBytes)
				log.Println(w.Header())
				return
			}
			responsebody, code := getuserstatus(userID[0])
			responsebodyBytes, _ := json.Marshal(responsebody)

			w.WriteHeader(code)
			w.Header().Set("Content-Type", "application/json")
			w.Write(responsebodyBytes)

		default:
			http.Error(w, "Invalid request method.", http.StatusBadRequest)
		}
	}
}

func handle_AddCoin() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			userID, exist := r.Header["Userid"]
			if !exist {
				log.Println("not userID")

				w.WriteHeader(http.StatusNotFound)
				return
			}
			respBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			var body User_CoinData
			json.Unmarshal(respBody, &body)
			saveUserCoin(userID[0], body.Coins)
			w.WriteHeader(http.StatusOK)
		}
	}
}

func getuserstatus(userID string) (User_Status, int) {
	var status User_Status
	query := fmt.Sprintf("SELECT EXISTS (SELECT * FROM T_USERS_GAMEINFO WHERE id = '%s');", userID)
	row, err := db.SelectQueryRow(query)
	if err != nil {
		log.Println("db error : ", err)
		status.Error = err.Error()
		return status, http.StatusInternalServerError
	}
	var result int
	row.Scan(&result)
	if result == 0 {
		log.Println("User is not.")
		status.Error = "Please Sign up."
		return status, http.StatusNotImplemented
	}
	query = fmt.Sprintf("SELECT nick, status, room_id, coins FROM t_users_gameinfo WHERE id = '%s';", userID)
	row, err = db.SelectQueryRow(query)
	if err != nil {
		log.Println("db error : ", err)
		status.Error = err.Error()
		return status, http.StatusInternalServerError
	}
	var roomid sql.NullString
	row.Scan(&status.Nick, &status.Status, &roomid, &status.Coins)
	if status.Status == "not login" {
		query = fmt.Sprintf("UPDATE t_users_gameinfo SET status = 'wait' WHERE id = '%s';", userID)
		db.ExecuteQuery(query)
	} else if status.Status == "play" || status.Status == "room" {
		if roomid.Valid {
			status.Room = roomid.String
		} else {
			status.Room = ""
			query = fmt.Sprintf("UPDATE t_users_gameinfo SET status = 'wait' WHERE id = '%s';", userID)
			db.ExecuteQuery(query)
		}
	}
	return status, http.StatusOK
}

func saveUser(userID string, userNick string) int {
	query := fmt.Sprintf("INSERT INTO T_USERS_GAMEINFO (`id`, `nick`, `coins`) VALUES ('%s', '%s', '30');", userID, userNick)
	result, err := db.ExecuteQuery(query)
	if err != nil {
		//error
		log.Println("Error : ", err)
		return http.StatusInternalServerError
	}
	if result == 0 {
		log.Println("Error : Check query. query : ", query)
	}

	return http.StatusOK
}

func saveUserCoin(userID string, coins int) int {
	query := fmt.Sprintf("SELECT coins FROM T_USERS_GAMEINFO WHERE id = '%s'", userID)
	row, _ := db.SelectQueryRow(query)
	var user_coin int
	row.Scan(&user_coin)

	user_coin += coins

	query = fmt.Sprintf("UPDATE T_USERS_GAMEINFO SET coins = '%d' WHERE id = '%s'", user_coin, userID)
	db.ExecuteQuery(query)

	return http.StatusOK
}
