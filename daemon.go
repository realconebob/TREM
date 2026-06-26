package main

// daemon.go - Background service that interacts with trem. Reads the reminder file and executes a reminder when the time is reached

import (
	"bytes"
	"encoding/gob"
	"os"
	"time"
)

type TremDaemon struct {
	Reminders []ReminderEntry
	CurrentTime time.Time
	ReminderFile *WatchedFile
}

func CreateDaemonFromBuffer(data []byte) (TremDaemon, error) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	var reminders []ReminderEntry

	if err := decoder.Decode(&reminders); err != nil {
		return TremDaemon{}, err
	}

	return TremDaemon{
		Reminders: reminders,
		CurrentTime: time.Now(),
		ReminderFile: &WatchedFile{},
	}, nil
}

func CreateDaemonFromFile(path string, poll time.Duration) (TremDaemon, error) {
	contents, err := os.ReadFile(path)
	if err != nil {return TremDaemon{}, err}

	res, err := CreateDaemonFromBuffer(contents)
	if err != nil {return TremDaemon{}, err}

	file, err := GetFileWatch(path, poll)
	if err != nil {return TremDaemon{}, err} // I could leave out this check by returning err, but I won't

	res.ReminderFile = file
	return res, nil
}

func (d *TremDaemon) Run() {
	// The daemon should:
		// *Get* populated with the current reminder file's reminders (This is not done by the daemon itself)
		// Fire the reminders at their designated times, or immediately if the reminder time has passed and TriggerIfMissed is true
		// Notice updates to the reminders file and update its list of reminders accordingly (without duplicating reminders!)
		// Save its list of reminders to the reminders file before shutting down
}