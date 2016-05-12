package main

import(
	"log"
	"github.com/henderjon/shutdown"
	"time"
)

func main(){

	signal := make(shutdown.SignalChan)

	go func(){
		for {
			log.Println("Zzzz")
			time.Sleep(time.Second * 5)
		}
	}()

	shutdown.Watch(signal, func() {
		log.Println("3...2...1...")
	})

}
