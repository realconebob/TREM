package main

import (
	"github.com/realconebob/trem/internal/cli"

	"fmt"
	"os"
)

const DEFAULT_REMINDER_FILE string = "./REMINDERFILE.test"
const DEFAULT_DAEMON_FILE string = "./DAEMONFILE.test"

func main() {
	res := cli.ProcCLIArgs()
	if res.Err != nil {
		fmt.Fprintf(os.Stderr, "Encountered an error: %v\n", res.Err)
		os.Exit(1)
	}

	var err error
	switch res.Command {
	case cli.TC_ADD:	err = cli.AddReminder(DEFAULT_REMINDER_FILE, res.Arguments)
	case cli.TC_EDIT:	err = cli.EditReminder(DEFAULT_REMINDER_FILE, res.Arguments)
	case cli.TC_LIST:	err = cli.ListReminders(DEFAULT_REMINDER_FILE, res.Arguments)
	case cli.TC_REMOVE:	err = cli.RemoveReminder(DEFAULT_REMINDER_FILE, res.Arguments)
	case cli.TC_DAEMON:	err = cli.CommandDaemon(DEFAULT_DAEMON_FILE, res.Arguments)
	default:			err = res.Err
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Encountered an error: %v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
