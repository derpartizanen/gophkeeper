package main

import (
	"os"

	"github.com/derpartizanen/gophkeeper/internal/keeperctl/controller/cmdline"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
)

func main() {
	if err := cmdline.Execute(buildVersion, buildDate); err != nil {
		os.Exit(1)
	}
}
