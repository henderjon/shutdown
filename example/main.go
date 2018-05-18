package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/henderjon/shutdown"
)

func destruct() {
	// this is the destructor before shutting down
	fmt.Printf("3...2...1... ")
	time.Sleep(time.Second * 2)
	fmt.Printf("and done.\n")
}

func main() {

	var wg sync.WaitGroup
	shutdown := shutdown.New(destruct)

	for n := 0; n < 13; n++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("Zzzz")
			time.Sleep(time.Second * 3)
		}()
	}

	wg.Wait()
	shutdown.Now("test example is over")
}
