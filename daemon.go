package main

// daemon.go - Background service that interacts with trem. Reads the reminder file and executes a reminder when the time is reached

import (
	"bytes"
	"encoding/gob"
	"os"
	"time"
)


type ReminderEntry struct {
	Registered time.Time 			`gob:"str"`
	TriggerOn time.Time				`gob:"end"`
	Message string					`gob:"msg"`
	Triggered bool					`gob:"tgr"`
	TriggerIfMissed bool			`gob:"mis"`

	// TODO: Unsupported as of now. Also unoptimal packing/memory usage if left here
	RepeatInterval time.Duration	`gob:"-"`
	Repeat bool						`gob:"-"`
}

func CreateReminderEntry(triggerOn time.Time, message string) ReminderEntry {
	return ReminderEntry{
		Registered: time.Now(),
		TriggerOn: triggerOn,
		Message: message,
		Triggered: false,
		TriggerIfMissed: true,

		Repeat: false,
		RepeatInterval: 0,
	}
}

func CreateReminderEntry_ByDate(layout, value, message string) (ReminderEntry, error) {
	parsed, err := time.Parse(layout, value)
	if err != nil {return ReminderEntry{}, err}
	return CreateReminderEntry(parsed, message), nil
}

func SerializeRemindersToFile(reminders []ReminderEntry, path string) error {
	file, err := os.Create(path)
	if err != nil {return err}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(reminders); err != nil {return err}
	return nil
}


type TremDaemon struct {
	Reminders []ReminderEntry
	CurrentTime time.Time
	PollingPeriod time.Duration
	ReminderFile string
}

func CreateDaemonFromBuffer(data []byte, polling time.Duration) (TremDaemon, error) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	var reminders []ReminderEntry

	if err := decoder.Decode(&reminders); err != nil {
		return TremDaemon{}, err
	}

	return TremDaemon{
		Reminders: reminders,
		CurrentTime: time.Now(),
		PollingPeriod: polling,
		ReminderFile: "",
	}, nil
}

func CreateDaemonFromFile(path string, polling time.Duration) (TremDaemon, error) {
	contents, err := os.ReadFile(path)
	if err != nil {return TremDaemon{}, err}

	// I could check for an error here, but there's not much point when the person running this function should check the error return value over the struct for indication of an error
	res, err := CreateDaemonFromBuffer(contents, polling)
	res.ReminderFile = path
	return res, err
}