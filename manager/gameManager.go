package manager

import (
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func Handler_Game_Ready() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
		}
	}
}

func Handler_Game_Status() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:

		}
	}
}

func Handler_Game_Endturn() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:

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
	timeSource := rand.NewSource(time.Now().UnixNano())
	random := rand.New(timeSource)
	card := shuffle(random)
	dbCard := strings.Join(card, "-")
	newCard := strings.Split(dbCard, "-")
	log.Print(newCard)
}
