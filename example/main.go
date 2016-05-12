package main

import(
	"log"
	"github.com/henderjon/shutdown"
	"time"
	"sync"
)

func main(){

	var wg sync.WaitGroup
	signal := make(shutdown.SignalChan)

	wg.Add(1)
	go func(signal shutdown.SignalChan){
		for {
			select {
			case <-signal:
				time.Sleep(time.Second * 3)
				wg.Done()
			default:
				log.Println("Zzzz")
				time.Sleep(time.Second * 3)
			}
		}
	}(signal)

	shutdown.Watch(signal, func() {
		log.Println("3...2...1...")
		wg.Wait()
		log.Println("and done.")
	})

}
