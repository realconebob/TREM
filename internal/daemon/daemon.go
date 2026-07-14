package daemon

// daemon.go - Background service that interacts with trem. Reads the reminder file and executes a reminder when the time is reached

import (
	"github.com/realconebob/trem/internal/reminders"
	"github.com/realconebob/trem/internal"
	"errors"
	"time"
)

type TremDaemon struct {
	Reminders []reminders.Entry
	Dispatched map[uint64]reminders.Entry
	CurrentTime time.Time
	CommandFile *misc.WatchedFile
	ReminderFile string
}

func CreateDaemonFromBuffer(data []byte) (TremDaemon, error) {
	rems, err := reminders.GetFromGob(data)

	return TremDaemon{
		Reminders: rems,
		Dispatched: make(map[uint64]reminders.Entry, 0),
		CurrentTime: time.Now(),
		CommandFile: &misc.WatchedFile{},
		ReminderFile: "",
	}, err
}

func CreateDaemonFromFile(cfile, rfile string, poll time.Duration) (TremDaemon, error) {
	rems, err := reminders.GetFromGobFile(rfile)
	if err != nil {return TremDaemon{}, err}

	watch, err := misc.GetFileWatch(cfile, poll)
	if err != nil {return TremDaemon{}, err}

	return TremDaemon{
		Reminders: rems,
		Dispatched: make(map[uint64]reminders.Entry, 0),
		CurrentTime: time.Now(),
		CommandFile: watch,
		ReminderFile: rfile,
	}, nil
}

func (d *TremDaemon) DispatchReminders() {
	if d == nil {return}
	for _, reminder := range d.Reminders {
		if _, ok := d.Dispatched[reminder.Identifier]; ok {continue} // Much better
		d.Dispatched[reminder.Identifier] = reminder

		go func(cr *reminders.Entry){
			time.Sleep(time.Until(cr.TriggerOn))
			cr.Triggered = true
			// TODO: Emit a notification or pop up an application or window with the reminder
		}(&reminder)
	}
}

func (d *TremDaemon) UpdateReminders() error {
	if d == nil {return errors.New("d is nil")}
	// Update daemon with ReminderEntry-s from file
	newReminders, err := reminders.GetFromGobFile(d.ReminderFile)
	if err != nil {return err}
	d.Reminders = newReminders

	return nil
}

func (d *TremDaemon) Close() error {
	if d == nil {return errors.New("d is nil")}
	d.CommandFile.Close()
	return reminders.SerializeToGobFile(d.Reminders, d.ReminderFile)
}

func (d *TremDaemon) UpdateCommands() error {
	if d == nil {return errors.New("d is nil")}
	return errors.New("<TremDaemon::UpdateCommands> Error: Not Implemented")
	// TODO: Implement
}
func (d *TremDaemon) RunCommands() error {
	if d == nil {return errors.New("d is nil")}
	return errors.New("<TremDaemon::RunCommands> Error: Not Implemented")
	// TODO: Implement
}

func (d *TremDaemon) Run() {
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
		<- fileUpdate
		// TODO: Add error checking for these
		d.UpdateCommands()
		d.RunCommands()
	}
}