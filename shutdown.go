package shutdown

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

const prefixSignal = "signal:"

// Shutdown listens for SIGINT and SIGTERM and executes the Destructor
type Shutdown struct {
	Signal   chan bool
	Destruct func()
	exit     func(int)
	*log.Logger
}

// New generates a new Shutdown with typical defaults
func New(destruct func()) *Shutdown {
	down := &Shutdown{
		Signal:   make(chan bool),
		Destruct: destruct,
		Logger:   log.New(os.Stderr, "", log.LUTC|log.LstdFlags),
		exit:     os.Exit, // if we embed this, we can mock it in our test #WINNING
	}
	go down.listen()
	return down
}

// Listen watches for os.Interrupt (syscall.SIGINT) and os.Kill (syscall.SIGTERM)
// [doc](https://golang.org/pkg/os/#Signal).
func (shutdown *Shutdown) listen() {

	sysSigChan := make(chan os.Signal)
	signal.Notify(sysSigChan, syscall.SIGINT)
	signal.Notify(sysSigChan, syscall.SIGTERM)

	select {
	// block for a signal
	case sig := <-sysSigChan:
		shutdown.Now(sig.String())
	// block until the application calls Now()
	case <-shutdown.Signal:
	}
}

// Now allows an application to trigger it's own shutdown.
func (shutdown *Shutdown) Now(reason string) {
	close(shutdown.Signal)
	shutdown.Destruct()
	shutdown.Println(prefixSignal, reason)
	shutdown.exit(1)
}
