package manager

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

type endTurn_requestBody struct {
	RoomID string `json:"room_id"`
	Speak  string `json:"speak"`
	MyBet  int    `json:"my_bet"`
	MyCoin int    `json:"my_coin"`
}
type game_requestBody struct {
	RoomID string `json:"room_id"`
}

type status_responseBody struct {
	Status      string `json:"status"`
	User        string `json:"user"`
	MyStatus    string `json:"my_status"`
	EnemyStatus string `json:"enemy_status"`
	EnemyBet    int    `json:"enemy_bet"`
	EnemyCoin   int    `json:"enemy_coin"`
	Error       struct {
		Code   int    `json:"code"`
		String string `json:"string"`
	} `json:"error"`
}

type result_responseBody struct {
	Result  string `json:"result"`
	GetCoin int    `json:"get_coin"`
	Coin    int    `json:"coin"`
	Error   struct {
		Code   int    `json:"code"`
		String string `json:"string"`
	} `json:"error"`
}

func Handler_Game_Ready() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			w.WriteHeader(http.StatusOK)
		}
	}
}

func Handler_Game_Status() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			var body game_requestBody
			respBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println(err)
			}
			json.Unmarshal(respBody, &body)

			w.WriteHeader(http.StatusOK)
		}
	}
}

func Handler_Game_Endturn() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			var body endTurn_requestBody
			respBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println(err)
			}
			json.Unmarshal(respBody, &body)
			w.WriteHeader(http.StatusOK)
		}
	}
}

func Handler_Game_Result() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			var body game_requestBody
			respBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println(err)
			}
			json.Unmarshal(respBody, &body)
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

func forgame() {
	//카드생성
	timeSource := rand.NewSource(time.Now().UnixNano())
	random := rand.New(timeSource)
	card := shuffle(random)
	dbCard := strings.Join(card, "-")
	newCard := strings.Split(dbCard, "-")
	log.Print(newCard)

	//방이름
	room_id, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
	}
	log.Println(room_id)
}
