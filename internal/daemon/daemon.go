package daemon

// daemon.go - Background service that interacts with trem. Reads the reminder file and executes a reminder when the time is reached

import (
	"errors"
	"fmt"
	"time"

	"github.com/realconebob/trem/internal"
	"github.com/realconebob/trem/internal/reminders"
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
	Behavior CommandIdx
	Path string
}

type Config struct {
	ReminderFile string
	ConfigFile string
}

// Note: This struct is specifically for dealing with daemon state. Things not related to state belong in the config struct
type Daemon struct {
	Settings Config
	Reminders []reminders.Entry
	Dispatched map[uint64]reminders.Entry
	CommandFile *misc.WatchedFile
	paused bool
	shutdown bool
}

func CreateDaemonFromBuffer(data []byte) (Daemon, error) {
	rems, err := reminders.GetFromGob(data)

	return Daemon{
		Reminders: rems,
		Dispatched: make(map[uint64]reminders.Entry, 0),
		CommandFile: &misc.WatchedFile{},
	}, err
}

func CreateDaemonFromFile(cfile, rfile string, poll time.Duration) (Daemon, error) {
	rems, err := reminders.GetFromGobFile(rfile)
	if err != nil {return Daemon{}, err}

	watch, err := misc.GetFileWatch(cfile, poll)
	if err != nil {return Daemon{}, err}

	return Daemon{
		Settings: Config{
			ReminderFile: rfile,
		},
		Reminders: rems,
		Dispatched: make(map[uint64]reminders.Entry, 0),
		CommandFile: watch,
	}, nil
}

func (d *Daemon) InterpretCommand(cmd Command) error {
	if d == nil {return errors.New("d is nil")}
	switch cmd.Behavior {
	case RELOAD_REMINDERS: 	return d.ReloadReminders()
	case RELOAD_CONFIG:		return d.ReloadConfig(cmd.Path)
	case SAVE_REMINDERS:	return d.SaveReminders(cmd.Path)
	case PAUSE:				return d.Pause()
	case SHUTDOWN:			return d.Close()
	case UNSPEC:			return errors.New("Unspecified command")
	case UNKNOWN: fallthrough
	default: return errors.New("Unknown command: " + fmt.Sprint(cmd.Behavior))
	}
}

func (d *Daemon) DispatchReminders() {
	if d == nil {return}
	for _, reminder := range d.Reminders {
		if _, ok := d.Dispatched[reminder.Identifier]; ok {continue} // Much better
		d.Dispatched[reminder.Identifier] = reminder

		go func(cr *reminders.Entry){
			if d.paused {return} // Don't bother if the daemon is paused

			time.Sleep(time.Until(cr.TriggerOn))
			cr.Triggered = true
			// TODO: Emit a notification or pop up an application or window with the reminder
		}(&reminder)
	}
}

func (d *Daemon) ReloadReminders() error {
	if d == nil {return errors.New("d is nil")}
	// Update daemon with ReminderEntry-s from file
	newReminders, err := reminders.GetFromGobFile(d.Settings.ReminderFile)
	if err != nil {return err}
	d.Reminders = newReminders

	return nil
}

func (d *Daemon) Close() error {
	if d == nil {return errors.New("d is nil")}
	if d.shutdown == true {return errors.New("d has already been shutdown/closed")}

	d.shutdown = true
	d.CommandFile.Close()
	return reminders.SerializeToGobFile(d.Reminders, d.Settings.ReminderFile)
}

func (d *Daemon) UpdateCommands() error {
	if d == nil {return errors.New("d is nil")}
	return errors.New("<Daemon::UpdateCommands> Error: Not Implemented")
	// TODO: Implement
}
func (d *Daemon) RunCommands() error {
	if d == nil {return errors.New("d is nil")}
	return errors.New("<Daemon::RunCommands> Error: Not Implemented")
	// TODO: Implement
}

func (d *Daemon) ReloadConfig(path string) error {
	// TODO:
	panic("Unimplemented")
}
func (d *Daemon) SaveReminders(path string) error {
	if path == "" {path = d.Settings.ReminderFile}
	return reminders.SerializeToGobFile(d.Reminders, path)
}
func (d *Daemon) Pause() error {
	d.paused = !d.paused
	return nil
}

func (d *Daemon) Run() {
	if d == nil {panic("d is nil")}
	// The daemon should:
		// *Get* populated with the current reminder file's reminders (This is not done by the daemon itself)
		// Fire the reminders at their designated times, or immediately if the reminder time has passed and TriggerIfMissed is true
		// Notice updates to the reminders file and update its list of reminders accordingly (without duplicating reminders!)
		// Save its list of reminders to the reminders file before shutting down

	defer d.Close()

	// Start watching the file from a goroutine
	fileUpdate := make(chan bool)
	go func(){
		res, err := d.CommandFile.CheckForUpdate()
		if err != nil {panic("Could not check reminder file for updates")}
		if res {fileUpdate <- true}
		d.CommandFile.Sleep()
	}()

	for {
		if d.shutdown {break}
		<- fileUpdate
		// TODO: Add error checking for these
		d.UpdateCommands()
		d.RunCommands()
	}
}