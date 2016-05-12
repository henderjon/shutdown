package main

import (
	"fmt"
	"github.com/henderjon/shutdown"
	"sync"
	"time"
)

func main() {

	var wg sync.WaitGroup
	signal := make(shutdown.SignalChan)

	wg.Add(1)
	go func(signal shutdown.SignalChan) {
		for {
			select {
			case <-signal:
				time.Sleep(time.Second * 3)
				wg.Done()
			default:
				fmt.Println("Zzzz")
				time.Sleep(time.Second * 3)
			}
		}
	}(signal)

	shutdown.Watch(signal, func() {
		fmt.Printf("3...2...1... ")
		wg.Wait()
		fmt.Printf("and done.\n")
	})

}
