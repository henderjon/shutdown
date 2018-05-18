package shutdown

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Shutdown listens for SIGINT and SIGTERM and executes the Destructor
type Shutdown struct {
	Signal   chan bool
	Destruct func()
	*log.Logger
}

// New generates a new Shutdown
func New(destruct func()) *Shutdown {
	down := &Shutdown{
		Signal:   make(chan bool),
		Destruct: destruct,
		Logger:   log.New(os.Stderr, "", log.LUTC|log.LstdFlags),
	}
	go down.listen()
	return down
}

// Listen is our signal watching func. It's worth noting that this is a
// blocking action so it should be run in a goroutine or as the last function call
// such that all the other goroutines will continue working while this func blocks.
// This func also watches for os.Interrupt (syscall.SIGINT) and os.Kill (syscall.SIGTERM)
// [doc](https://golang.org/pkg/os/#Signal).
func (shutdown *Shutdown) listen() {

	sysSigChan := make(chan os.Signal)
	signal.Notify(sysSigChan, syscall.SIGINT)
	signal.Notify(sysSigChan, syscall.SIGTERM)

	select {
	// block for a signal
	case sig := <-sysSigChan:
		shutdown.Now(sig.String())
	case <-shutdown.Signal:
	}
}

// Now allows an application to trigger it's own shutdown.
func (shutdown *Shutdown) Now(reason string) {
	close(shutdown.Signal)
	shutdown.Destruct()
	shutdown.Println("signal:", reason)
	os.Exit(1)
}
