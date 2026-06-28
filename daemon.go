package main

// daemon.go - Background service that interacts with trem. Reads the reminder file and executes a reminder when the time is reached

import (
	"time"
)

type TremDaemon struct {
	Reminders []ReminderEntry
	Dispatched map[uint64]ReminderEntry
	CurrentTime time.Time
	ReminderFile *WatchedFile
}

func CreateDaemonFromBuffer(data []byte) (TremDaemon, error) {
	reminders, err := GetRemindersFromGob(data)

	return TremDaemon{
		Reminders: reminders,
		Dispatched: make(map[uint64]ReminderEntry, 0),
		CurrentTime: time.Now(),
		ReminderFile: &WatchedFile{},
	}, err
}

func CreateDaemonFromFile(path string, poll time.Duration) (TremDaemon, error) {
	reminders, err := GetRemindersFromGobFile(path)
	if err != nil {return TremDaemon{}, err}

	file, err := GetFileWatch(path, poll)
	if err != nil {return TremDaemon{}, err}

	return TremDaemon{
		Reminders: reminders,
		Dispatched: make(map[uint64]ReminderEntry, 0),
		CurrentTime: time.Now(),
		ReminderFile: file,
	}, nil
}

func (d *TremDaemon) DispatchReminders() {
	for _, reminder := range d.Reminders {
		if _, ok := d.Dispatched[reminder.Identifier]; ok {continue} // Much better
		d.Dispatched[reminder.Identifier] = reminder

		go func(cr *ReminderEntry){
			time.Sleep(time.Until(cr.TriggerOn))
			cr.Triggered = true
			// TODO: Emit a notification or pop up an application or window with the reminder
		}(&reminder)
	}
}

func (d *TremDaemon) UpdateReminders() error {
	// Get file contents to make buffer
	file := d.ReminderFile.Handle
	stat, err := file.Stat()
	if err != nil {return err}

	// Populate buffer
	_, err = file.Seek(0, 0)
	if err != nil {return err}
	buf := make([]byte, stat.Size())
	if _, err := file.Read(buf); err != nil {return err}

	// Update daemon with ReminderEntry-s from buffer
	newReminders, err := GetRemindersFromGob(buf)
	if err != nil {return err}
	d.Reminders = newReminders

	return nil
}

func (d *TremDaemon) Shutdown() error {
	rmpath := d.ReminderFile.Handle.Name()
	d.ReminderFile.Close()
	return SerializeRemindersToGobFile(d.Reminders, rmpath)
}

func (d *TremDaemon) Run() {
	// The daemon should:
		// *Get* populated with the current reminder file's reminders (This is not done by the daemon itself)
		// Fire the reminders at their designated times, or immediately if the reminder time has passed and TriggerIfMissed is true
		// Notice updates to the reminders file and update its list of reminders accordingly (without duplicating reminders!)
		// Save its list of reminders to the reminders file before shutting down

	defer d.Shutdown()

	// Start watching the file from a goroutine
	fileUpdate := make(chan bool)
	go func(){
		res, err := d.ReminderFile.CheckForUpdate()
		if err != nil {panic("Could not check reminder file for updates")}
		if res {fileUpdate <- true}
		time.Sleep(d.ReminderFile.PollingInterval)
	}()

	// Wait for file updates, update reminders accordingly
	for {
		<- fileUpdate
		d.UpdateReminders()
		d.DispatchReminders()
	}
}