package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/henderjon/shutdown"
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
				// pause before returning, we're shutting down
				time.Sleep(time.Second * 3)
				return
			default:
				// this is where you do your work because the chan isn't closed yet
				fmt.Println("Zzzz")
				time.Sleep(time.Second * 3)
			}
			if n == 5 {
				// shutdown ater 5 loops
				shutdown.Now(signal)
				return
			}
		}
	}(signal)

	// if you're doing anything else in your application (e.g. a web server) you'll want to `go` here
	shutdown.Watch(signal, func() {
		// this is the destructor before shutting down
		fmt.Printf("3...2...1... ")
		wg.Wait()
		fmt.Printf("and done.\n")
	})

}
