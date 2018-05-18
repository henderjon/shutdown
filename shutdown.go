package shutdown

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const prefixSignal = "signal:"

// Shutdown listens for SIGINT and SIGTERM and executes the Destructor.
type Shutdown struct {
	Destruct func()
	signal   chan bool
	exit     func(int)
	once     sync.Once
	*log.Logger
}

// New generates a new Shutdown with typical defaults
func New(destruct func()) *Shutdown {
	down := &Shutdown{
		signal:   make(chan bool),
		Destruct: destruct,
		Logger:   log.New(os.Stderr, "", log.LUTC|log.LstdFlags),
		exit:     os.Exit, // if we embed this, we can mock it in our test #WINNING
	}
	go down.listen()
	return down
}

// Now allows an application to trigger it's own shutdown.
func (shutdown *Shutdown) Now(reason string) {
	shutdown.once.Do(func() {
		shutdown.now(reason)
	})
}

// IsDown checks to see if the shutdown has been triggered
func (shutdown *Shutdown) IsDown() bool {
	select {
	case <-shutdown.signal:
		return true
	default:
		return false
	}
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
		shutdown.once.Do(func() {
			shutdown.now(sig.String())
		})
	// block until the application calls Now()
	case <-shutdown.signal:
	}
}

// now wraps our shutdown logic in a sync.Once
func (shutdown *Shutdown) now(reason string) {
	close(shutdown.signal)
	shutdown.Destruct()
	shutdown.Println(prefixSignal, reason)
	shutdown.exit(1)
}
