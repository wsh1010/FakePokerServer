package main

import (
	"log"

	"github.com/gofrs/uuid"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	defer func() {
		if r := recover(); r != nil {
			log.Println("recover : ", r)

			main()
		}
	}()

	room_id, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
	}
	log.Println(room_id)
}
