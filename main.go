package main

import (
	"log"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	defer func() {
		if r := recover(); r != nil {
			log.Println("recover : ", r)

			main()
		}
	}()

}
