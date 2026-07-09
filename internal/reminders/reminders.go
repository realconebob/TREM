package reminders

// reminders.go - Representation of a text reminder

import (
	"github.com/realconebob/trem/internal"
	"encoding/gob"
	"bytes"
	"time"
	"fmt"
	"os"
)

type Entry struct {
	Identifier      uint64    `gob:"id"`
	Registered      time.Time `gob:"str"`
	TriggerOn       time.Time `gob:"end"`
	Message         string    `gob:"msg"`
	Triggered       bool      `gob:"tgr"`
	TriggerIfMissed bool      `gob:"mis"`

	// TODO: Unsupported as of now. Also unoptimal packing/memory usage if left here
	RepeatInterval time.Duration `gob:"-"`
	Repeat         bool          `gob:"-"`
}

func CreateEntry(triggerOn time.Time, message string) Entry {
	now := time.Now()
	return Entry{
		Identifier:      misc.PseudoRandom(uint64(now.UnixNano())),
		Registered:      now,
		TriggerOn:       triggerOn,
		Message:         message,
		Triggered:       false,
		TriggerIfMissed: true,

		Repeat:         false,
		RepeatInterval: 0,
	}
}

func CreateEntryByDate(layout, value, message string) (Entry, error) {
	parsed, err := time.Parse(layout, value)
	if err != nil {
		return Entry{}, err
	}
	return CreateEntry(parsed, message), nil
}

func (entryA Entry) Compare(entryB Entry) bool {
	return (entryA.Identifier == entryB.Identifier) && (entryA.TriggerOn.Equal(entryB.TriggerOn))
}

func (entry Entry) String() string {
	return fmt.Sprintf(
		"ReminderFile@%p{id: %v, registered: %v, triggeron: %v, message: %v, triggered: %v, tmissed: %v}",
		&entry, entry.Identifier, entry.Registered, entry.TriggerOn, entry.Message, entry.Triggered, entry.TriggerIfMissed,
	)
}

func SerializeToGobFile(reminders []Entry, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	reminders = misc.Filter(reminders, func(entry Entry) bool {
		// Only include entries that have yet to be triggered
		return entry.Triggered == false
	})
	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(reminders); err != nil {
		return err
	}
	return nil
}

func GetFromGob(data []byte) ([]Entry, error) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	var reminders []Entry

	if err := decoder.Decode(&reminders); err != nil {
		return []Entry{}, err
	}

	return reminders, nil
}

func GetFromGobFile(path string) ([]Entry, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return []Entry{}, err
	}
	return GetFromGob(contents)
}
