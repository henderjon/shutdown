package shutdown

import (
	"bytes"
	"log"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestShutdown(t *testing.T) {

	tables := []struct {
		val1, val2 string
	}{
		{"foo", "bar"},
		{"buzz", "bazz"},
	}

	for _, table := range tables {
		var destructVal string
		var logVal bytes.Buffer
		var count int

		shutdown := &Shutdown{
			signal: make(chan bool),
			Destruct: func() {
				destructVal = table.val1
			},
			exit: func(i int) {
				// do nothing
			},
			Logger: log.New(&logVal, "", 0),
		}

		for x := 0; x < 4; x++ {
			if !shutdown.IsDown() {
				count++
			}
		}

		shutdown.Now(table.val2)

		if diff := cmp.Diff(destructVal, table.val1); diff != "" {
			t.Errorf("shutdown.Destruct: (-got +want)\n%s", diff)
		}

		if diff := cmp.Diff(count, 4); diff != "" {
			t.Errorf("shutdown.Wait: (-got +want)\n%s", diff)
		}

		expected := prefixSignal + " " + table.val2 + "\n"
		if diff := cmp.Diff(logVal.String(), expected); diff != "" {
			t.Errorf("shutdown.Logger: (-got +want)\n%s", diff)
		}
	}
}
