package main

import (
	"os"

	"github.com/realconebob/trem/internal"
	"github.com/realconebob/trem/internal/cli"
	"github.com/realconebob/trem/internal/daemon"
)



func main() {
	conf, err := daemon.GetDefaultConfig()
	misc.PrintErrAndExit(err, "Could not get default config: %v", err)

	res := cli.ProcCLIArgs()
	misc.PrintErrAndExit(res.Err, "Encountered an error: %v\n", res.Err)

	err = res.Command(misc.Ternary(res.IsDaemonCmd, conf.CommandPath, conf.ReminderPath), res.Arguments)
	misc.PrintErrAndExit(err, "Could not process command: %v\n", err)

	os.Exit(0)
}
