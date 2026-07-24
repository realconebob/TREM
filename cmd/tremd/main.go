package main

import (
	"os"
	"github.com/realconebob/trem/internal/daemon"
)

func main() {
	// Note: Daemon should process cli arguments for alternate config locations. For now it's hard coded to be the user's default config

	os.Exit(daemon.Launch(os.Args))
}
