package main

import (
	"fmt"
	"os"
	"time"

	misc "github.com/realconebob/trem/internal"
	"github.com/realconebob/trem/internal/daemon"
)

func main() {
	// Note: Daemon should process cli arguments for alternate config locations. For now it's hard coded to be the user's default config

	conf, err := daemon.GetConfigFromFile(daemon.CONFIG_CONFIG_NAME)
	if err != nil {
		errs := daemon.GenerateDefaultConfig()
		errs = misc.NilErrSliceCheck(errs)
		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "Encountered error while generating default config: %v", e)
		}
		if errs != nil {os.Exit(1)}

		conf, err = daemon.GetDefaultConfig()
		if err != nil {panic(fmt.Sprintf("Could not get config: %v", err))}
	}

	daemon, err := daemon.NewFromConfig(conf, time.Second * 5)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Encountered error while creating daemon: %v", err)
		os.Exit(1)
	}

	daemon.Run()
	// This is deadlocking apparently. Wonderful
}
