package shutdown

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const prefixSignal = "signal:"

type Destructor func()

type Listener func(chan os.Signal)

// Shutdown listens for SIGINT and SIGTERM and executes the Destructor
type Shutdown struct {
	Destructor Destructor
	signal     chan bool
	block      chan bool
	exit       func(int)
	once       sync.Once
	*log.Logger
}

// New generates a new Shutdown with typical defaults
func New(destruct Destructor, signals []syscall.Signal) *Shutdown {
	down := &Shutdown{
		signal:     make(chan bool),
		Destructor: destruct,
		Logger:     log.New(os.Stderr, "\n", log.LUTC|log.LstdFlags),
		exit:       os.Exit, // if we embed this, we can mock it in our test #WINNING
		block:      make(chan bool),
	}

	go down.listen(signals)
	return down
}

// Now allows an application to trigger it's own shutdown
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

// Wait is a func that allows the calling context to block until shutdown is finished
func (shutdown *Shutdown) Wait() {
	<-shutdown.block
}

func (shutdown *Shutdown) SetDestructor(destruct Destructor) {
	shutdown.Destructor = destruct
}

// Listen watches for SIGINT SIGTERM [doc](https://golang.org/pkg/os/#Signal)
func (shutdown *Shutdown) listen(signals []syscall.Signal) {

	sysSigChan := make(chan os.Signal, 1)
	signal.Notify(sysSigChan, syscall.SIGINT)
	signal.Notify(sysSigChan, syscall.SIGTERM)
	for _, sig := range signals {
		signal.Notify(sysSigChan, sig)
	}

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
	shutdown.Println(prefixSignal, reason)
	close(shutdown.signal)
	if shutdown.Destructor != nil {
		shutdown.Destructor()
	}
	close(shutdown.block)
	shutdown.exit(1)
}

// wraps any request and checks to make sure that the server isn't shutting down
func (shutdown *Shutdown) Handler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if shutdown.IsDown() {
			http.Error(w, "The service is unavailable or shutting down", http.StatusServiceUnavailable)
		} else {
			handler.ServeHTTP(w, r)
		}
	})
}
