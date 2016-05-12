package shutdown

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	sysSigChan     chan os.Signal
	shutdownLogger = log.New(os.Stderr, "", 0) // log to stderr without the timestamps
)

// Destructor is a func takes no args and returns no values. It is executed as an injectable destructor
// so that the calling context could use this func to sync.Wait() or print a pretty exit message.
type Destructor func()

// SignalChan is a blank channel used to signal a shutdown
type SignalChan chan struct{}

// Watch is our signal watching func. It's worth noting that this is a
// blocking action so it should be run in a goroutine or as the last function call
// such that all the other goroutines can continue working while the application
// sits blocked here.
// For a deeper discussion of the close channel idiom: http://dave.cheney.net/2013/04/30/curious-channels
func Watch(shutdown SignalChan, destruct Destructor) {

	sysSigChan = make(chan os.Signal, 1)
	signal.Notify(sysSigChan, syscall.SIGINT)
	signal.Notify(sysSigChan, syscall.SIGTERM)

	var sig os.Signal

	select { // block until we get a signal
	case sig = <-sysSigChan:
		close(shutdown) // idiom via: http://dave.cheney.net/2013/04/30/curious-channels
	}

	shutdownLogger.Printf("\n.signal: %s; shutting down...\n", sig.String())
	destruct()
	shutdownLogger.Printf(".shutdown: program exit at %s\n", time.Now().Format(time.RFC3339))
	os.Exit(1)
}