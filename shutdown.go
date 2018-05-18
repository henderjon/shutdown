package shutdown

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Shutdown listens for SIGINT and SIGTERM and executes the Destructor
type Shutdown struct {
	toggle   chan bool
	destruct func()
}

// New generates a new Shutdown
func New(destruct func()) *Shutdown {
	down := &Shutdown{
		toggle:   make(chan bool),
		destruct: destruct,
	}
	go down.listen()
	return down
}

// listen is our signal watching func. It's worth noting that this is a
// blocking action so it should be run in a goroutine or as the last function call
// such that all the other goroutines will continue working while this func blocks.
// This func also watches for os.Interrupt (syscall.SIGINT) and os.Kill (syscall.SIGTERM)
// [doc](https://golang.org/pkg/os/#Signal).
func (shutdown *Shutdown) listen() {

	sysSigChan := make(chan os.Signal)
	signal.Notify(sysSigChan, syscall.SIGINT)
	signal.Notify(sysSigChan, syscall.SIGTERM)

	// block for a signal
	sig := <-sysSigChan
	shutdown.Now(sig.String())
}

// Now allows an application to trigger it's own shutdown.
func (shutdown *Shutdown) Now(reason string) {
	close(shutdown.toggle)
	shutdown.destruct()
	log.Printf("sig: %s; datetime: %s\n", reason, time.Now().UTC().Format(time.RFC3339))
	os.Exit(1)
}

// For a deeper discussion of the close channel idiom: http://dave.cheney.net/2013/04/30/curious-channels
