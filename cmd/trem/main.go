package main

import (
	"github.com/realconebob/trem/internal/cli"
	"github.com/realconebob/trem/internal/reminders"

	"errors"
	"fmt"
	"os"
	"time"
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
	case cli.TC_ADD:
		err = AddReminder(DEFAULT_REMINDER_FILE, res.Arguments)
	case cli.TC_EDIT:
		err = EditReminder(DEFAULT_REMINDER_FILE, res.Arguments)
	case cli.TC_LIST:
		err = ListReminders(DEFAULT_REMINDER_FILE, res.Arguments)
	case cli.TC_REMOVE:
		err = RemoveReminder(DEFAULT_REMINDER_FILE, res.Arguments)
	case cli.TC_DAEMON:
		err = CommandDaemon(DEFAULT_DAEMON_FILE, res.Arguments)
	default:
		err = res.Err
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Encountered an error: %v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func AddReminder(file string, args []string) error {
	currentReminders, err := reminders.GetFromGobFile(file)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	layout := args[0]
	switch layout {
	case "Layout":
		layout = time.Layout
	case "ANSIC":
		layout = time.ANSIC
	case "UnixDate":
		layout = time.UnixDate
	case "RubyDate":
		layout = time.RubyDate
	case "RFC822":
		layout = time.RFC822
	case "RFC822Z":
		layout = time.RFC822Z
	case "RFC850":
		layout = time.RFC850
	case "RFC1123":
		layout = time.RFC1123
	case "RFC1123Z":
		layout = time.RFC1123Z
	case "RFC3339":
		layout = time.RFC3339
	case "RFC3339Nano":
		layout = time.RFC3339Nano
	case "Kitchen":
		layout = time.Kitchen

	case "Stamp":
		layout = time.Stamp
	case "StampMilli":
		layout = time.StampMilli
	case "StampMicro":
		layout = time.StampMicro
	case "StampNano":
		layout = time.StampNano
	case "DateTime":
		layout = time.DateTime
	case "DateOnly":
		layout = time.DateOnly
	case "TimeOnly":
		layout = time.TimeOnly
	} // Pulled straight from https://pkg.go.dev/time#Layout

	// TODO: Make this more robust
	datevalue := args[1]
	message := args[2]
	newReminder, err := reminders.CreateEntryByDate(layout, datevalue, message)
	if err != nil {
		return err
	}

	currentReminders = append(currentReminders, newReminder)
	return reminders.SerializeToGobFile(currentReminders, file)
}
func EditReminder(file string, args []string) error {
	currentReminders, err := reminders.GetFromGobFile(file)
	if err != nil {
		return err
	}
	currentReminders = currentReminders // Shut go up for the time being
	// TODO: Implement

	return nil
}
func ListReminders(file string, args []string) error {
	reminders, err := reminders.GetFromGobFile(file)
	if err != nil {
		return err
	}

	for _, reminder := range reminders {
		fmt.Println(reminder)
	}

	return nil
}
func RemoveReminder(file string, args []string) error {
	currentReminders, err := reminders.GetFromGobFile(file)
	if err != nil {
		return err
	}
	currentReminders = currentReminders // Shut go up for the time being
	// TODO: Implement

	return nil
}

func CommandDaemon(file string, args []string) error {
	return errors.New("TODO: main::CommandDaemon is unimplemented")
}
