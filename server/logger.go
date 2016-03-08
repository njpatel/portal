package server

import (
	"os"

	"github.com/njpatel/loggo"
	"golang.org/x/crypto/ssh/terminal"
)

var logger = loggo.GetLogger("portal")

func init() {
	loggerSpec := os.Getenv("PORTAL_DEBUG")

	// Setup logging and such things if we're running in a term
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		if loggerSpec == "" {
			loggerSpec = "<root>=DEBUG"
		}
		// As we're in a terminal, let's make the output a little nicer
		_, _ = loggo.ReplaceDefaultWriter(loggo.NewSimpleWriter(os.Stderr, &loggo.ColorFormatter{}))
	} else {
		if loggerSpec == "" {
			loggerSpec = "<root>=WARNING"
		}
	}

	_ = loggo.ConfigureLoggers(loggerSpec)
}
