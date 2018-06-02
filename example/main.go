package main

import (
	"fmt"
	"log"
	"os"
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

	log.Println(os.Getpid())

	shutdown := shutdown.New(destruct)

	for n := 0; n < 13; n++ {
		go func() {
			fmt.Println("Zzzz")
			time.Sleep(time.Second * 3)
		}()
	}

	go func() {
		time.Sleep(time.Duration(4) * time.Second)
		shutdown.Now("the example is over")
	}()

	shutdown.Wait()

}
