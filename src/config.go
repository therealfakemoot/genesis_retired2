package genesis

import (
	//"github.com/spf13/cobra"
	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"os"
)

func SetupLogging() {
	log.SetHandler(text.New(os.Stderr))
}
