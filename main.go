package main

import (
	"log"
	"math"
	"math/rand"
	"strings"
	"time"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	defer func() {
		if r := recover(); r != nil {
			log.Println("recover : ", r)

			main()
		}
	}()
	compare()
}

func shuffle(r *rand.Rand) string {
	cards := make([]int, 0)
	for i := 0; i < 20; i++ {
		cardnum := i/2 + 1
		cards = append(cards, cardnum)
	}

	arr := []string{"1", "2", "3"}

	for i := len(arr) - 1; i > 0; i-- {
		j := int(math.Floor(r.Float64() * float64(i+1)))

		temp := arr[i]
		arr[i] = arr[j]
		arr[j] = temp
	}
	result := strings.Join(arr, "")
	return result
}

func compare() {
	results := make(map[string]int)
	timeSource := rand.NewSource(time.Now().UnixNano())
	random := rand.New(timeSource)
	for i := 0; i < 1000000; i++ {
		result := shuffle(random)
		if _, exist := results[result]; !exist {
			results[result] = 1
		} else {
			results[result]++
		}
	}
	for key, value := range results {
		log.Println(key, " : ", value)
	}
}
