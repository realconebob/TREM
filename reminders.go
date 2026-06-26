package main

// reminders.go - Representation of a text reminder

import (
	"time"
	"os"
	"encoding/gob"
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

	reminders = Filter(reminders, func(entry ReminderEntry)bool{
		// Only include entries that have yet to be triggered
		return entry.Triggered == false
	})
	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(reminders); err != nil {return err}
	return nil
}