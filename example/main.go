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
		defer wg.Done()
		for n := 0; ; n++ {
			select {
			case <-signal:
				time.Sleep(time.Second * 3)
				return
			default:
				fmt.Println("Zzzz")
				time.Sleep(time.Second * 3)
			}
			if n == 5 {
				shutdown.Now(signal)
				return
			}
		}
	}(signal)

	shutdown.Watch(signal, func() {
		fmt.Printf("3...2...1... ")
		wg.Wait()
		fmt.Printf("and done.\n")
	})

}
