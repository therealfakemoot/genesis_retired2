package log

import (
	"github.com/sirupsen/logrus"
	"os"
)

// Term is the terminal logger for the Genesis application.
var Term *logrus.Logger

func setupTerm() {
	Term = logrus.New()
	Term.Out = os.Stdout

}

func init() {
	setupTerm()
}
