package cli

// cli.go - functions and structs dealing with the command line tool for interacting with tremd
// TODO: Finish the daemon then revisit this mess

import (
	"github.com/realconebob/trem/internal/reminders"

	"errors"
	"time"
	"fmt"
	"os"
)

type TremCommand uint32

const (
	TC_UNSPEC TremCommand = iota
	TC_UNKNOWN

	TC_LIST
	TC_ADD
	TC_EDIT
	TC_REMOVE
	TC_DAEMON

	TC_TOOBIG
)

type CLIRes struct {
	Command   TremCommand
	Arguments []string
	Err       error
}

func CreateCLIResFromArgs(command string, args []string) CLIRes {
	var res CLIRes = CLIRes{Arguments: args}
	switch command {
	case "add":		res.Command = TC_ADD
	case "edit":	res.Command = TC_EDIT
	case "list":	res.Command = TC_LIST
	case "remove":	res.Command = TC_REMOVE
	case "daemon":	res.Command = TC_DAEMON
	default:
		res.Err = errors.New("Unknown command \"" + command + "\"")
		res.Command = TC_UNKNOWN
	}

	return res
}

func ProcCLIArgs() CLIRes {
	fmt.Println("Arguments:")
	for _, arg := range os.Args {
		fmt.Printf("\"%v\" ", arg)
	}
	fmt.Print("\n")

	if len(os.Args) < 2 {
		return CLIRes{Err: errors.New("Too few CLI arguments")}
	}
	workingArgs := os.Args[1:]
	return CreateCLIResFromArgs(workingArgs[0], workingArgs[1:])
}

func AddReminder(file string, args []string) error {
	currentReminders, err := reminders.GetFromGobFile(file)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	layout := args[0]
	switch layout {
	case "Layout":		layout = time.Layout
	case "ANSIC":		layout = time.ANSIC
	case "UnixDate":	layout = time.UnixDate
	case "RubyDate":	layout = time.RubyDate
	case "RFC822":		layout = time.RFC822
	case "RFC822Z":		layout = time.RFC822Z
	case "RFC850":		layout = time.RFC850
	case "RFC1123":		layout = time.RFC1123
	case "RFC1123Z":	layout = time.RFC1123Z
	case "RFC3339":		layout = time.RFC3339
	case "RFC3339Nano":	layout = time.RFC3339Nano
	case "Kitchen":		layout = time.Kitchen

	case "Stamp":		layout = time.Stamp
	case "StampMilli":	layout = time.StampMilli
	case "StampMicro":	layout = time.StampMicro
	case "StampNano":	layout = time.StampNano
	case "DateTime":	layout = time.DateTime
	case "DateOnly":	layout = time.DateOnly
	case "TimeOnly":	layout = time.TimeOnly
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
