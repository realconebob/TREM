package cli

// cli.go - functions and structs dealing with the command line tool for interacting with tremd
// TODO: Finish the daemon then revisit this mess
// TODO TODO: See above. It's getting worse

import (
	"github.com/realconebob/trem/internal/reminders"
	"github.com/realconebob/trem/internal"
	"strconv"
	"errors"
	"time"
	"fmt"
	"os"
)

type CLIRes struct {
	Command   func(string,[]string)error
	Arguments []string
	Err       error
}

func CreateCLIResFromArgs(command string, args []string) CLIRes {
	var res CLIRes = CLIRes{Arguments: args}
	switch command {
	case "add":		res.Command = AddReminder
	case "edit":	res.Command = EditReminder
	case "list":	res.Command = ListReminders
	case "del":		res.Command = RemoveReminder
	case "daemon":	res.Command = CommandDaemon
	default:
		res.Err = errors.New("Unknown command \"" + command + "\"")
		res.Command = func(_ string, _ []string)error {return errors.New("Got unknown command from CreateCLIResFromArgs")}
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

// Converts a time.Layout name to its respective literal string. Returns empty string if match not found
// Ex: "UnixDate" -> "Mon Jan _2 15:04:05 MST 2006"
func LayoutNameToLayoutLiteral(name string) string {
	var res string
	switch name {
	case "Layout":		res = time.Layout
	case "ANSIC":		res = time.ANSIC
	case "UnixDate":	res = time.UnixDate
	case "RubyDate":	res = time.RubyDate
	case "RFC822":		res = time.RFC822
	case "RFC822Z":		res = time.RFC822Z
	case "RFC850":		res = time.RFC850
	case "RFC1123":		res = time.RFC1123
	case "RFC1123Z":	res = time.RFC1123Z
	case "RFC3339":		res = time.RFC3339
	case "RFC3339Nano":	res = time.RFC3339Nano
	case "Kitchen":		res = time.Kitchen

	case "Stamp":		res = time.Stamp
	case "StampMilli":	res = time.StampMilli
	case "StampMicro":	res = time.StampMicro
	case "StampNano":	res = time.StampNano
	case "DateTime":	res = time.DateTime
	case "DateOnly":	res = time.DateOnly
	case "TimeOnly":	res = time.TimeOnly
	} // Pulled straight from https://pkg.go.dev/time#Layout

	return res
}

func AddReminder(file string, args []string) error {
	if l := len(args); l != 3 {return errors.New("Not enough arguments given. Expected: 3, Got: " + fmt.Sprint(l))}

	currentReminders, err := reminders.GetFromGobFile(file)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	layout := args[0]
	if converted := LayoutNameToLayoutLiteral(layout); len(converted) > 0 {layout = converted}

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

func IndexReminder(rems []reminders.Entry, index uint64) (*reminders.Entry, error) {
	var toEdit *reminders.Entry = nil
	if index >= uint64(len(rems)) {
		for _, rem := range rems {
			if rem.Identifier == index {toEdit = &rem}
		}
	} else {toEdit = &rems[index]}
	if toEdit == nil {return nil, errors.New("Could not figure out which reminder to edit. Index: " + fmt.Sprint(index))}

	return toEdit, nil
}

func EditReminder(file string, args []string) error {
	// All operations require at least 3 arguments. The "time" operation needs 4
	arglen := len(args); if arglen < 3 {return errors.New("Not enough arguments given")}

	currentReminders, err := reminders.GetFromGobFile(file)
	if err != nil {
		return err
	}

	index, err := strconv.ParseUint(args[0], 10, 0)
	if err != nil {return err}

	toEdit, err := IndexReminder(currentReminders, index)
	if err != nil {return err}

	// determine edit operation:
	// valid ops:
		// Update time to trigger
		// Update message
		// Update trigger-if-missed

	switch args[1] {
	case "time":
		// use args[2] and args[3] to parse a new time value and update the trigger on time
		if arglen < 4 {return errors.New("Not enough arguments given")}
		layout, value := args[2], args[3]
		if converted := LayoutNameToLayoutLiteral(layout); len(converted) > 0 {layout = converted}

		newTrigger, err := time.Parse(layout, value)
		if err != nil {return err}
		toEdit.TriggerOn = newTrigger

	case "msg":
		// use args[2] to update the message
		toEdit.Message = args[2]

	case "miss":
		// use args[2] to update the trigger-if-missed flag
		flag, err := strconv.ParseBool(args[2])
		if err != nil {return err}
		toEdit.TriggerIfMissed = flag

	default: return errors.New("Unrecognized edit directive: \"" + args[1] + "\"")
	}

	return reminders.SerializeToGobFile(currentReminders, file)
}

func ListReminders(file string, args []string) error {
	reminders, err := reminders.GetFromGobFile(file)
	if err != nil {
		return err
	}

	// TODO: there should be optional arguments where you can change how everything is listed
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

	index, err := strconv.ParseUint(args[0], 10, 0)
	if err != nil {return err}

	toDel, err := IndexReminder(currentReminders, index)
	if err != nil {return err}

	currentReminders = misc.Filter(currentReminders, func(entry reminders.Entry)bool {
		return !toDel.Compare(entry)
	})

	return reminders.SerializeToGobFile(currentReminders, file)
}

func CommandDaemon(file string, args []string) error {
	return errors.New("TODO: main::CommandDaemon is unimplemented")
}
