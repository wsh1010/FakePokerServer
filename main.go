package main

import (
	"fakepokerserver/manager"
	"fakepokerserver/module/db"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	defer func() {
		if r := recover(); r != nil {
			log.Println("recover : ", r)
			time.Sleep(time.Second * 1)
			main()
		}
	}()

	done := make(chan int, 1)
	done <- 1
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(2)
	db.InitDB()
	go manager.RunningServer(&wg, done)
	go manager.CheckDB(&wg, done)

	<-signals
	<-done
	wg.Wait()
}
