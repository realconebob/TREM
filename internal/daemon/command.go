package daemon

import (
	"github.com/realconebob/trem/internal/reminders"
	"github.com/realconebob/trem/internal"
	"errors"
	"fmt"
)

type CommandIdx uint8
const (
	UNSPEC CommandIdx = iota
	RELOAD_REMINDERS
	RELOAD_CONFIG
	SAVE_REMINDERS
	PAUSE
	SHUTDOWN
	UNKNOWN
)
type Command struct {
	Behavior 	CommandIdx	`gob:"b"`
	Path 		string		`gob:"p"`
	executed 	bool		`gob:"e"`
}

func (d *Daemon) InterpretCommand(cmd *Command) error {
	if d == nil {return errors.New("d is nil")}
	cmd.executed = true
	switch cmd.Behavior {
	case RELOAD_REMINDERS: 	return d.ReloadReminders()
	case RELOAD_CONFIG:		return d.ReloadConfig(cmd.Path)
	case SAVE_REMINDERS:	return d.SaveReminders(cmd.Path)
	case PAUSE:				return d.Pause()
	case SHUTDOWN:			return d.Close()
	case UNSPEC:			return errors.New("Unspecified command")

	case UNKNOWN: 	fallthrough
	default: 		return errors.New("Unknown command: " + fmt.Sprint(cmd.Behavior))
	}
}

// Really there should only be a single command in the command file queue at a time, but this can process more if need be
func (d *Daemon) RunCommands() []error {
	if d == nil {return []error{errors.New("d is nil")}}
	// Read commands from command file, then interpret them
	commands, err := reminders.GetFromGobFile[Command](d.Settings.CommandPath)
	if err != nil {return []error{err}}

	errs := make([]error, 0)
	for _, command := range commands {
		errs = append(errs, d.InterpretCommand(&command))
	}

	// Clear the file of executed commands
	err = reminders.SerializeToGobFile(commands, d.Settings.CommandPath, func(cmd Command)bool {return !cmd.executed})
	errs = append(errs, err)

	errs = misc.Filter(errs, func(err error)bool {return err != nil})
	if len(errs) <= 0 {return nil}
	return errs
}