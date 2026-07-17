package reminders

// reminders.go - Representation of a text reminder

import (
	"github.com/realconebob/trem/internal"
	"time"
	"fmt"
)

type Entry struct {
	Identifier      uint64    `gob:"id"`
	Registered      time.Time `gob:"str"`
	TriggerOn       time.Time `gob:"end"`
	Message         string    `gob:"msg"`
	Triggered       bool      `gob:"tgr"`
	TriggerIfMissed bool      `gob:"mis"`

	// TODO: Unsupported as of now. Also unoptimal packing/memory usage if left here
	RepeatInterval	time.Duration 	`gob:"-"`
	RepeatToCount	uint			`gob:"-"`
	RepeatCount		uint			`gob:"-"`
	Repeat			bool          	`gob:"-"`
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

// Get the function used to filter reminder entries when serializing to a gob
func GetReminderGobFilter() func(Entry)bool {
	// Only include entries that have yet to be triggered
	return func(entry Entry) bool {return entry.Triggered == false}
}
