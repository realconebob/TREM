package main

// daemon.go - Background service that interacts with trem. Reads the reminder file and executes a reminder when the time is reached

import (
	"time"
)

type TremDaemon struct {
	Reminders []ReminderEntry
	Dispatched map[ReminderEntry]bool
	CurrentTime time.Time
	ReminderFile *WatchedFile
}

func CreateDaemonFromBuffer(data []byte) (TremDaemon, error) {
	reminders, err := GetRemindersFromGob(data)

	return TremDaemon{
		Reminders: reminders,
		Dispatched: make(map[ReminderEntry]bool, 0),
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
		Dispatched: make(map[ReminderEntry]bool, 0),
		CurrentTime: time.Now(),
		ReminderFile: file,
	}, nil
}

func (d *TremDaemon) DispatchReminders() {
	for _, reminder := range d.Reminders {
		// This is terrible but I don't know of a better solution rn
		var skip bool
		for entry := range d.Dispatched {
			if entry.Compare(reminder) {
				skip = true
				break // Note: Works as expected, only breaks this inner loop
			}
		}
		if skip {continue}

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
// Man, go's zero values gave me an awesome idea that I can't do within go (or maybe any language for that) matter. The idea is an
// "Updatable" wrapper for any type. The wrapper would store any generic type and implement said generic type's signature through
// delegation. Then the Updatable class would contain some flags and behaviors for accepting or rejecting updates to the held
// generic. Say a "set once" flag that lets the value be set, but then further attempts to update it are ignored. That way, in
// combination with go's nice zero values, I could basically just write the same function but remove all the "if err != nil"
// statements littered in it and return the first error that came up. Alas, I must clutter my code

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

	fileUpdate := make(chan bool)
	// Start watching the file
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