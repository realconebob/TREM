package daemon

// daemon.go - Background service that interacts with trem. Reads the reminder file and executes a reminder when the time is reached

import (
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/realconebob/trem/internal"
	"github.com/realconebob/trem/internal/gobwrap"
	"github.com/realconebob/trem/internal/reminders"
)

var defaultConf Config
func init() {
	var err error
	defaultConf, err = GetDefaultConfig()
	if err != nil {fmt.Fprintf(os.Stderr, "Could not get path to default config locations: %v", err)}
}

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

type Config struct {
	ReminderPath string
	CommandPath string
	ConfigPath string
}

func GetDefaultConfig() (Config, error) {
	confdir, err := os.UserConfigDir()
	if err != nil {return Config{}, err}

	rems := path.Join(confdir, "trem", "reminders.gob")
	cmds := path.Join(confdir, "trem", "commands.gob")
	conf := path.Join(confdir, "trem", "tremrc")

	return Config{
		ReminderPath: rems,
		CommandPath: cmds,
		ConfigPath: conf,
	}, nil
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

	dfCopy := defaultConf
	return Daemon{
		Settings: dfCopy,
		Reminders: rems,
		Dispatched: make(map[uint64]reminders.Entry, 0),
		CommandFile: &misc.WatchedFile{},
	}, err
}

func CreateDaemonFromFile(cmdfile, rfile string, poll time.Duration) (Daemon, error) {
	rems, err := reminders.GetFromGobFile(rfile)
	if err != nil {return Daemon{}, err}

	watch, err := misc.GetFileWatch(cmdfile, poll)
	if err != nil {return Daemon{}, err}

	dfCopy := defaultConf
	dfCopy.ReminderPath = rfile
	return Daemon{
		Settings: dfCopy,
		Reminders: rems,
		Dispatched: make(map[uint64]reminders.Entry, 0),
		CommandFile: watch,
	}, nil
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
	newReminders, err := reminders.GetFromGobFile(d.Settings.ReminderPath)
	if err != nil {return err}
	d.Reminders = newReminders

	return nil
}

func (d *Daemon) Close() error {
	if d == nil {return errors.New("d is nil")}
	if d.shutdown == true {return errors.New("d has already been shutdown/closed")}

	d.shutdown = true
	d.CommandFile.Close()
	return reminders.SerializeToGobFile(d.Reminders, d.Settings.ReminderPath)
}

// Really there should only be a single command in the command file queue at a time, but this can process more if need be
func (d *Daemon) RunCommands() []error {
	if d == nil {return []error{errors.New("d is nil")}}
	// Read commands from command file, then interpret them
	commands, err := gobwrap.GetFromGobFile[Command](d.Settings.CommandPath)
	if err != nil {return []error{err}}

	errs := make([]error, 0)
	for _, command := range commands {
		errs = append(errs, d.InterpretCommand(&command))
	}

	// Clear the file of executed commands
	err = gobwrap.SerializeToGobFile(commands, d.Settings.CommandPath, func(cmd Command)bool {return !cmd.executed})
	errs = append(errs, err)

	errs = misc.Filter(errs, func(err error)bool {return err != nil})
	if len(errs) <= 0 {return nil}
	return errs
}

func (d *Daemon) ReloadConfig(path string) error {
	if d == nil {return errors.New("d is nil")}

	// TODO:
	panic("Unimplemented")
}
func (d *Daemon) SaveReminders(path string) error {
	if d == nil {return errors.New("d is nil")}
	if path == "" {path = d.Settings.ReminderPath}
	return reminders.SerializeToGobFile(d.Reminders, path)
}
func (d *Daemon) Pause() error {
	if d == nil {return errors.New("d is nil")}
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
		// TODO: Add better error checking
		errs := d.RunCommands()
		for _, err := range errs {
			fmt.Fprintf(os.Stderr, "Warning - Could not process a command: %v", err)
		}
	}
}