package main

import (
	"fmt"
	"os"
)

const DEFAULT_REMINDER_FILE string = "./REMINDER_FILE.test"

func main() {
	res := ProcCLIArgs()
	if res.err != nil {fmt.Fprintf(os.Stderr, "Encountered an error: %v\n", res.err); os.Exit(1)}

	var err error
	switch res.command {
	case TC_ADD: 	err = AddReminder(DEFAULT_REMINDER_FILE, res.arguments)
	case TC_EDIT: 	err = EditReminder(DEFAULT_REMINDER_FILE, res.arguments)
	case TC_LIST: 	err = ListReminders(DEFAULT_REMINDER_FILE, res.arguments)
	case TC_REMOVE: err = RemoveReminder(DEFAULT_REMINDER_FILE, res.arguments)
	default: 		err = res.err
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Encountered an error: %v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func AddReminder(file string, args []string) error {
	currentReminders, err := GetRemindersFromGobFile(file)
	if err != nil {return err}
	currentReminders = currentReminders // Shut go up for the time being
	// TODO: Implement

	return nil
}
func EditReminder(file string, args []string) error {
	currentReminders, err := GetRemindersFromGobFile(file)
	if err != nil {return err}
	currentReminders = currentReminders // Shut go up for the time being
	// TODO: Implement

	return nil
}
func ListReminders(file string, args []string) error {
	reminders, err := GetRemindersFromGobFile(file)
	if err != nil {return err}

	for reminder := range reminders {
		fmt.Println(reminder)
	}

	return nil
}
func RemoveReminder(file string, args []string) error {
	currentReminders, err := GetRemindersFromGobFile(file)
	if err != nil {return err}
	currentReminders = currentReminders // Shut go up for the time being
	// TODO: Implement

	return nil
}