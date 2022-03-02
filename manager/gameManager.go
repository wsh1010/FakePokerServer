package manager

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"fakepokerserver/module/db"

	"github.com/gofrs/uuid"
)

type status_user struct {
	UserID     string
	UserNick   string
	UserStatus string
	UserBet    int
	Usercoins  int
}

type endTurn_requestBody struct {
	RoomID string `json:"room_id"`
	Speak  string `json:"speak"`
	MyBet  int    `json:"bet_coin"`
}
type game_requestBody struct {
	RoomID string `json:"room_id"`
	Status string `json:"status"`
}

type status_responseBody struct {
	Status      string `json:"status"`
	Mynick      string `json:"my_nick"`
	MyStatus    string `json:"my_status"`
	MyBet       int    `json:"my_bet"`
	MyCoins     int    `json:"my_coins"`
	MyCard      int    `json:"my_card"`
	Enemynick   string `json:"enemy_nick"`
	EnemyStatus string `json:"enemy_status"`
	EnemyBet    int    `json:"enemy_bet"`
	EnemyCoins  int    `json:"enemy_coins"`
	EnemyCard   int    `json:"enemy_card"`
	Error       string `json:"error"`
}

type roomResult_responseBdy struct {
	RoomID string `json:"room_id"`
	Error  string `json:"error"`
}

type result_responseBody struct {
	Result  string `json:"result"`
	GetCoin int    `json:"get_coins"`
	Coins   int    `json:"coins"`
	Error   string `json:"error"`
}

func Handler_Game_Ready() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			var body game_requestBody
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
			result := SetUserRoomStatus(body.RoomID, userID[0], body.Status, 1)
			w.WriteHeader(result)
		}
	}
}

func Handler_Game_Status() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userID, exist := r.Header["Userid"]
			if !exist {
				log.Println("not userID")
				var requestBody status_responseBody
				requestBody.Error = "not userID"
				requestBodyBytes, _ := json.Marshal(requestBody)
				w.WriteHeader(http.StatusNotImplemented)
				w.Write(requestBodyBytes)
				return
			}
			URISplit := strings.Split(r.RequestURI, "/")
			roomID := URISplit[len(URISplit)-1]
			responseBody, result := GetRoomStatus(roomID, userID[0])

			responseBodyBytes, _ := json.Marshal(responseBody)
			w.WriteHeader(result)
			w.Write(responseBodyBytes)
			//log.Println(string(responseBodyBytes))
		}
	}
}

func Handler_Game_Endturn() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			var body endTurn_requestBody
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
			ProcessResult(userID[0], body)
			w.WriteHeader(http.StatusOK)

		}
	}
}

func Handler_Game_Result() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userID, exist := r.Header["Userid"]
			if !exist {
				log.Println("not userID")
				w.WriteHeader(http.StatusNotImplemented)
				return
			}
			URISplit := strings.Split(r.RequestURI, "/")
			roomID := URISplit[len(URISplit)-1]

			result := SetUserRoomInit(roomID, userID[0], "wait", "room")
			w.WriteHeader(result)
		}
	}
}

func Handle_joinRoom() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var responseBody roomResult_responseBdy
		switch r.Method {
		case http.MethodPut:
			userID, exist := r.Header["Userid"]
			if !exist {
				log.Println("not userID")
				responseBody.Error = "not userID"
				responseBodyBytes, _ := json.Marshal(responseBody)
				w.WriteHeader(http.StatusNotImplemented)
				w.Write(responseBodyBytes)
				return
			}
			responseBody = JoinRoom(userID[0])
			responseBodyBytes, _ := json.Marshal(responseBody)
			w.WriteHeader(http.StatusOK)
			w.Write(responseBodyBytes)
		case http.MethodPost:
			userID, exist := r.Header["Userid"]
			if !exist {
				log.Println("not userID")
				responseBody.Error = "not userID"
				responseBodyBytes, _ := json.Marshal(responseBody)
				w.WriteHeader(http.StatusNotImplemented)
				w.Write(responseBodyBytes)
				return
			}
			responseBody = CreateRoom(userID[0])
			responseBodyBytes, _ := json.Marshal(responseBody)
			w.WriteHeader(http.StatusOK)
			w.Write(responseBodyBytes)
		}
	}
}

func Handle_ExitRoom() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			userID, exist := r.Header["Userid"]
			if !exist {
				log.Println("not userID")
				w.WriteHeader(http.StatusNotImplemented)
				return
			}
			URISplit := strings.Split(r.RequestURI, "/")
			roomID := URISplit[len(URISplit)-1]
			ExitRoom(roomID, userID[0])
			w.WriteHeader(http.StatusOK)
		}
	}
}

func shuffle(r *rand.Rand) []string {
	cards := make([]string, 0)
	for i := 0; i < 20; i++ {
		card_num := i/2 + 1
		cards = append(cards, strconv.Itoa(card_num))
	}

	for i := len(cards) - 1; i > 0; i-- {
		j := int(math.Floor(r.Float64() * float64(i+1)))

		temp := cards[i]
		cards[i] = cards[j]
		cards[j] = temp
	}
	return cards
}

func makeCards() string {
	//카드생성
	timeSource := rand.NewSource(time.Now().UnixNano())
	random := rand.New(timeSource)
	card := shuffle(random)
	dbCard := strings.Join(card, "-")
	return dbCard
	//방이름
	/*room_id, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
	}
	log.Println(room_id)*/
}

func CreateRoom(userID string) roomResult_responseBdy {
	var result roomResult_responseBdy
	query := fmt.Sprintf("SELECT status FROM t_users_gameinfo where id = '%s';", userID)
	row, err := db.SelectQueryRow(query)
	if err != nil {
		log.Println("error : ", err)
		result.Error = err.Error()
		return result
	}
	var status string
	row.Scan(&status)
	if status != "wait" {
		log.Println("User is not create Room. " + status)
		result.Error = "User is not create Room."
		return result
	}

	room_id, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
	}

	dbCard := makeCards()
	query = fmt.Sprintf("INSERT INTO t_rooms_info VALUES ('%s', '%s', NULL, '%d', '%s', NULL);", room_id.String(), "user_wait", 0, dbCard)
	db.ExecuteQuery(query)
	query = fmt.Sprintf("UPDATE t_users_gameinfo SET room_id = '%s', status = 'room', room_status = 'wait' WHERE id = '%s'", room_id.String(), userID)
	db.ExecuteQuery(query)
	result.RoomID = room_id.String()
	return result
}

func JoinRoom(userID string) roomResult_responseBdy {
	var result roomResult_responseBdy
	query := "SELECT room_id FROM t_rooms_info WHERE status = 'user_wait';"
	rows, err := db.SelectQueryRows(query)
	if err != nil {
		log.Println(err)
	}
	var rooms []string
	for rows.Next() {
		var room string
		rows.Scan(&room)
		rooms = append(rooms, room)
	}
	size := len(rooms)
	if size == 0 {
		return CreateRoom(userID)
	} else {
		timeSource := rand.NewSource(time.Now().UnixNano())
		random := rand.New(timeSource)
		room_num := random.Intn(size)

		query = fmt.Sprintf("UPDATE t_users_gameinfo SET room_id = '%s', status = 'room', room_status = 'wait' WHERE id = '%s'", rooms[room_num], userID)
		db.ExecuteQuery(query)
		query = fmt.Sprintf("SELECT COUNT(*) from t_users_gameinfo where room_id = '%s'", rooms[room_num])
		row, _ := db.SelectQueryRow(query)
		var user_num int
		row.Scan(&user_num)
		if user_num == 2 {
			query := fmt.Sprintf("UPDATE t_rooms_info SET status = 'ready_wait' where room_id = '%s';", rooms[room_num])
			db.ExecuteQuery(query)
		}
		result.RoomID = rooms[room_num]
		return result
	}
}

func ExitRoom(roomID string, userID string) {
	query := fmt.Sprintf("SELECT id, status, bet, coins FROM t_users_gameinfo WHERE room_id = '%s'", roomID)
	rows, _ := db.SelectQueryRows(query)
	var me status_user
	var you status_user
	for rows.Next() {
		var id string
		var status string
		var bet int
		var coin int
		rows.Scan(&id, &status, &bet, &coin)

		if id == userID {
			me.UserID = id
			me.UserStatus = status
			me.UserBet = bet
			me.Usercoins = coin
		} else {
			you.UserID = id
			you.UserStatus = status
			you.UserBet = bet
			you.Usercoins = coin
		}
	}
	if me.UserStatus == "room" {
		me.Usercoins += me.UserBet
		me.UserBet = 0
		query = fmt.Sprintf("UPDATE t_users_gameinfo SET bet = '%d', coins = '%d' where room_id = '%s' and id = '%s'", me.UserBet, me.Usercoins, roomID, userID)
		db.ExecuteQuery(query)
	} else if me.UserStatus == "play" {
		me.UserStatus = "die"
		GameOver(roomID, me, you)
	}

	dbCard := makeCards()
	query = fmt.Sprintf("UPDATE t_rooms_info SET status = 'user_wait', round='0', cards='%s', start_player = NULL where room_id = '%s';", dbCard, roomID)
	db.ExecuteQuery(query)
	query = fmt.Sprintf("UPDATE t_users_gameinfo SET status = 'wait', room_id = null where room_id = '%s' and id = '%s'", roomID, userID)
	db.ExecuteQuery(query)
}

func StartGame(roomID string, users_id []string) {
	var start_Player sql.NullString
	query := fmt.Sprintf("SELECT start_player, round FROM t_rooms_info WHERE room_id = '%s';", roomID)
	row, _ := db.SelectQueryRow(query)
	var round int
	row.Scan(&start_Player, &round)

	if !start_Player.Valid {
		timeSource := rand.NewSource(time.Now().UnixNano())
		random := rand.New(timeSource)
		start_player_num := random.Intn(2)
		start_Player.String = users_id[start_player_num]
		round++
		start_Player.Valid = true
		query = fmt.Sprintf("UPDATE t_rooms_info SET status = 'play', start_player = '%s', round = '%d' WHERE room_id = '%s'", start_Player.String, round, roomID)
		db.ExecuteQuery(query)
	} else {
		if round == 10 {
			round = 1
			dbCards := makeCards()
			query = fmt.Sprintf("UPDATE t_rooms_info SET status = 'play', cards = '%s', round = '%d' WHERE room_id = '%s'", dbCards, round, roomID)
			db.ExecuteQuery(query)
		} else {
			round++
			query = fmt.Sprintf("UPDATE t_rooms_info SET status = 'play', round = '%d' WHERE room_id = '%s'", round, roomID)
			db.ExecuteQuery(query)
		}
	}

	for _, id := range users_id {
		if id == start_Player.String {
			query = fmt.Sprintf("UPDATE t_users_gameinfo SET status = 'play', room_status = 'turn' WHERE room_id = '%s' and id = '%s'", roomID, id)
		} else {
			query = fmt.Sprintf("UPDATE t_users_gameinfo SET status = 'play', room_status = 'wait' WHERE room_id = '%s' and id  = '%s'", roomID, id)
		}
		db.ExecuteQuery(query)
	}
}

func GetRoomStatus(roomID string, userID string) (status_responseBody, int) {
	var result status_responseBody
	var users_id []string
	exist_user := false
	rows, err := db.SelectQueryRows("SELECT t_users_gameinfo.`id`, t_users_gameinfo.`nick` ,t_rooms_info.`status`, t_users_gameinfo.`room_status` ,t_users_gameinfo.bet, t_users_gameinfo.coins, t_rooms_info.round, t_rooms_info.cards, t_rooms_info.start_player FROM t_rooms_info LEFT JOIN t_users_gameinfo ON t_rooms_info.room_id = t_users_gameinfo.room_id WHERE t_users_gameinfo.room_id = '" + roomID + "';")
	if err != nil {
		log.Println(err)
		var error_result status_responseBody
		error_result.Error = "DB Error"
		return error_result, http.StatusInternalServerError
	}
	defer rows.Close()

	for rows.Next() {
		var user status_user
		var game_status string
		var cards string
		var round int
		var start_player sql.NullString
		err = rows.Scan(&user.UserID, &user.UserNick, &game_status, &user.UserStatus, &user.UserBet, &user.Usercoins, &round, &cards, &start_player)
		Cards := strings.Split(cards, "-")
		cardnum := (round - 1) * 2
		if err != nil {
			log.Println(err)
			var error_result status_responseBody
			error_result.Error = "RoomID Not Found"
			return error_result, http.StatusNotImplemented
		}
		if user.UserID == userID {
			result.Mynick = user.UserNick
			result.MyStatus = user.UserStatus
			result.MyBet = user.UserBet
			result.MyCoins = user.Usercoins
			result.Status = game_status
			exist_user = true
			if game_status == "result" {
				if start_player.String != userID {
					cardnum += 1
				}
				result.MyCard, _ = strconv.Atoi(Cards[cardnum])
			} else {
				result.MyCard = 0
			}
		} else {
			result.Enemynick = user.UserNick
			result.EnemyBet = user.UserBet
			result.EnemyCoins = user.Usercoins
			result.EnemyStatus = user.UserStatus
			if game_status == "play" || game_status == "result" {
				if start_player.String != user.UserID {
					cardnum += 1
				}
				result.EnemyCard, _ = strconv.Atoi(Cards[cardnum])
			} else {
				result.EnemyCard = 0
			}
		}
		users_id = append(users_id, user.UserID)
	}

	if !exist_user {
		var error_result status_responseBody
		log.Println("user not found")
		error_result.Error = "User Not Found"
		return error_result, http.StatusNotImplemented
	}

	if result.EnemyStatus == "ready" && result.MyStatus == "ready" {
		// 둘다 레디면 게임 시작 DB로 변동
		StartGame(roomID, users_id)
	}

	return result, http.StatusOK
}

func SetUserRoomInit(roomID string, userID string, room_status string, status string) int {
	query := fmt.Sprintf("UPDATE t_users_gameinfo SET status = '%s', room_status = '%s', bet = '0' where id = '%s' and room_id = '%s';", status, room_status, userID, roomID)
	db.ExecuteQuery(query)
	query = fmt.Sprintf("UPDATE t_rooms_info SET status = 'ready_wait' WHERE room_id = '%s'", roomID)
	db.ExecuteQuery(query)

	return http.StatusOK
}

func SetUserRoomStatus(roomID string, userID string, status string, bet int) int {
	query := fmt.Sprintf("SELECT room_status, bet, coins from t_users_gameinfo where id = '%s' and room_id = '%s';", userID, roomID)
	row, err := db.SelectQueryRow(query)
	if err != nil {
		return http.StatusInternalServerError
	}
	var u_bet, u_coins int
	var u_status string
	row.Scan(&u_status, &u_bet, &u_coins)
	if u_status == status {
		log.Println("already")
		return http.StatusBadRequest
	}
	if status == "ready" {
		u_bet += bet
		u_coins -= bet
	} else {
		u_coins += u_bet
		u_bet = 0
	}

	query = fmt.Sprintf("UPDATE t_users_gameinfo SET bet = '%d', coins = '%d', room_status = '%s' where id = '%s' and room_id = '%s';", u_bet, u_coins, status, userID, roomID)
	success, err := db.ExecuteQuery(query)

	if err != nil {
		log.Println("db Error")
		return http.StatusInternalServerError
	}
	if success == 0 {
		log.Println("Not found id or room")
		return http.StatusNotImplemented
	}

	return http.StatusOK
}

func ProcessResult(userID string, result endTurn_requestBody) {
	var user_You status_user
	var user_I status_user
	query := fmt.Sprintf("SELECT id, bet, coins FROM t_users_gameinfo WHERE room_id='%s';", result.RoomID)
	rows, err := db.SelectQueryRows(query)
	if err != nil {
		return
	}
	for rows.Next() {
		var id string
		var bet int
		var coins int

		rows.Scan(&id, &bet, &coins)
		if id == userID {
			user_I.UserID = id
			user_I.UserBet = bet
			user_I.Usercoins = coins
		} else {
			user_You.UserID = id
			user_You.UserBet = bet
			user_You.Usercoins = coins
		}
	}
	defer rows.Close()
	user_I.UserBet += result.MyBet
	user_I.Usercoins -= result.MyBet

	if user_I.UserBet == user_You.UserBet || result.Speak == "die" {
		// 게임 종료
		user_I.UserStatus = result.Speak
		GameOver(result.RoomID, user_I, user_You)
		return
	}
	query = fmt.Sprintf("UPDATE t_users_gameinfo SET room_status = 'wait', bet = '%d', coins = '%d' where id = '%s' and room_id = '%s';",
		user_I.UserBet, user_I.Usercoins, user_I.UserID, result.RoomID)
	db.ExecuteQuery(query)
	query = fmt.Sprintf("UPDATE t_users_gameinfo SET room_status = 'turn' where id = '%s' and room_id = '%s';", user_You.UserID, result.RoomID)
	db.ExecuteQuery(query)
}

func GameOver(room_id string, p1 status_user, p2 status_user) {
	if p1.UserStatus == "die" {
		totalbets := p1.UserBet + p2.UserBet
		p1.UserBet = 0
		p1.UserStatus = "lose"
		p2.UserBet = totalbets
		p2.Usercoins += totalbets
		p2.UserStatus = "win"
	} else {
		row, err := db.SelectQueryRow("SELECT round, cards, start_player from t_rooms_info")
		if err != nil {
			log.Println(err)
		}
		var card_string string
		var round int
		var player string
		err = row.Scan(&round, &card_string, &player)
		if err != nil {
			log.Println(err)
		}
		Cards := strings.Split(card_string, "-")
		cardnum := (round - 1) * 2
		player1_card, _ := strconv.ParseInt(Cards[cardnum], 10, 64)
		player2_card, _ := strconv.ParseInt(Cards[cardnum+1], 10, 64)
		totalbets := p1.UserBet + p2.UserBet
		if player1_card > player2_card {
			if p1.UserID == player {
				//1플레이어 승
				p1.UserBet = totalbets
				p1.Usercoins += totalbets
				p1.UserStatus = "win"
				p2.UserBet = 0
				p2.UserStatus = "lose"
			} else {
				//2플레이어 승
				p1.UserBet = 0
				p1.UserStatus = "lose"
				p2.UserBet = totalbets
				p2.Usercoins += totalbets
				p2.UserStatus = "win"
			}

		} else if player1_card < player2_card {
			if p1.UserID == player {
				//2플레이어 승
				p1.UserBet = 0
				p1.UserStatus = "lose"
				p2.UserBet = totalbets
				p2.Usercoins += totalbets
				p2.UserStatus = "win"
			} else {
				//1플레이어 승
				p1.UserBet = totalbets
				p1.UserStatus = "win"
				p1.Usercoins += totalbets
				p2.UserBet = 0
				p2.UserStatus = "lose"
			}

		} else if player1_card == player2_card {
			// draw
			p1.UserStatus = "draw"
			p2.UserStatus = "draw"
		}
	}
	query := fmt.Sprintf("UPDATE t_users_gameinfo SET room_status = '%s', bet = '%d', coins = '%d' WHERE id = '%s' and room_id = '%s';",
		p1.UserStatus, p1.UserBet, p1.Usercoins, p1.UserID, room_id)
	db.ExecuteQuery(query)
	query = fmt.Sprintf("UPDATE t_users_gameinfo SET room_status = '%s', bet = '%d', coins = '%d' WHERE id = '%s' and room_id = '%s';",
		p2.UserStatus, p2.UserBet, p2.Usercoins, p2.UserID, room_id)
	db.ExecuteQuery(query)
	query = fmt.Sprintf("UPDATE t_rooms_info SET status = 'result' WHERE room_id = '%s';", room_id)
	db.ExecuteQuery(query)
}
