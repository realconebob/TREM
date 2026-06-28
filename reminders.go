package main

// reminders.go - Representation of a text reminder

import (
	"encoding/gob"
	"bytes"
	"time"
	"os"
)

type ReminderEntry struct {
	Identifier 		uint64		`gob:"id"`
	Registered		time.Time 	`gob:"str"`
	TriggerOn		time.Time	`gob:"end"`
	Message			string		`gob:"msg"`
	Triggered		bool		`gob:"tgr"`
	TriggerIfMissed	bool		`gob:"mis"`

	// TODO: Unsupported as of now. Also unoptimal packing/memory usage if left here
	RepeatInterval 	time.Duration	`gob:"-"`
	Repeat 			bool			`gob:"-"`
}

func CreateReminderEntry(triggerOn time.Time, message string) ReminderEntry {
	now := time.Now()
	return ReminderEntry{
		Identifier: PseudoRandom(uint64(now.UnixNano())),
		Registered: now,
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

func (entryA ReminderEntry) Compare(entryB ReminderEntry) bool {
	return (entryA.Identifier == entryB.Identifier) && (entryA.TriggerOn.Equal(entryB.TriggerOn))
}

func SerializeRemindersToGobFile(reminders []ReminderEntry, path string) error {
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

func GetRemindersFromGob(data []byte) ([]ReminderEntry, error) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	var reminders []ReminderEntry

	if err := decoder.Decode(&reminders); err != nil {
		return []ReminderEntry{}, err
	}

	return reminders, nil
}

func GetRemindersFromGobFile(path string) ([]ReminderEntry, error) {
	contents, err := os.ReadFile(path)
	if err != nil {return []ReminderEntry{}, err}
	return GetRemindersFromGob(contents)
}