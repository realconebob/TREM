package main

import (
	"github.com/realconebob/trem/internal/cli"
	"github.com/realconebob/trem/internal"
	"os"
)

const DEFAULT_REMINDER_FILE string = "./REMINDERFILE.test"
const DEFAULT_DAEMON_FILE string = "./DAEMONFILE.test"

func main() {
	res := cli.ProcCLIArgs()
	misc.PrintErrAndExit(res.Err, "Encountered an error: %v\n", res.Err)

	err := res.Command(DEFAULT_REMINDER_FILE, res.Arguments) // Note: Daemon commands not implemented yet, so passing in only the default reminder is fine, but this does need to be fixed in the future
	misc.PrintErrAndExit(err, "Could not process command: %v\n", err)

	os.Exit(0)
}
