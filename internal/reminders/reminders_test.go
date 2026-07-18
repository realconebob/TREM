package reminders

import (
	"github.com/realconebob/trem/internal"
	"testing"
	"path"
	"time"
	"os"
)

func Test_CreateEntry(t *testing.T) {
	ctime := time.Now()
	msg := "Sample message"
	res := CreateEntry(ctime, msg)

	if !res.TriggerOn.Equal(ctime) {
		t.Error("Reminder did not correctly store TriggerOn time")
	}
	if res.Message != msg {
		t.Error("Reminder did not correctly store user message")
	}
}

func Test_CreateEntryByDate(t *testing.T) {
	msg := "Sample msg"
	_, err := CreateEntryByDate(time.UnixDate, "Wed Jul 8 14:15:00 CST 2026", msg)
	if err != nil {
		t.Errorf("Reminder could not parse time it should be able to: %v", err)
	}

	if _, err := CreateEntryByDate(time.UnixDate, "7", msg); err == nil {
		t.Errorf("Reminder parsed a time it shouldn't have been able to: %v", err)
	}
}

func Test_ReminderEntry_Compare(t *testing.T) {
	e1 := CreateEntry(time.Now(), "entry 1")
	e2 := e1
	e3 := CreateEntry(time.Now().Add(time.Minute), "entry 3")

	if !e1.Compare(e1) {
		t.Errorf("Comparison between e1 and itself returned false when it should be true")
	}
	if !e1.Compare(e2) {
		t.Errorf("Comparison between e1 and e2 returned false when it should be true")
	}
	if e1.Compare(e3) {
		t.Errorf("Comparison between e1 and e3 returned true when it should be false")
	}

	if !e2.Compare(e2) {
		t.Errorf("Comparison between e2 and itself returned false when it should be true")
	}
	if !e2.Compare(e1) {
		t.Errorf("Comparison between e2 and e1 returned false when it should be true")
	}
	if e2.Compare(e3) {
		t.Errorf("Comparison between e2 and e3 returned true when it should be false")
	}

	if !e3.Compare(e3) {
		t.Errorf("Comparison between e3 and itself returned false when it should be true")
	}
	if e3.Compare(e1) {
		t.Errorf("Comparison between e3 and e1 returned true when it should be false")
	}
	if e3.Compare(e2) {
		t.Errorf("Comparison between e3 and e2 returned true when it should be false")
	}

}

func Test_Serialize(t *testing.T) {
	// SerializeToGobFile
	dir := t.TempDir()
	name := path.Join(dir, "reminders.gob")

	rems := []Entry{
		CreateEntry(time.Now().Add(time.Minute), "Entry 1"),
		CreateEntry(time.Now().Add(5 * time.Minute), "Entry 2"),
	}

	if err := misc.SerializeToGobFile(rems, name, GetReminderGobFilter()); err != nil {
		t.Errorf("Could not serialize reminders to gob file: %v", err)
	}


	// GetFromGobFile
	rems2, err := misc.GetFromGobFile[Entry](name)
	if err != nil {t.Errorf("Could not read reminders from gob: %v", err)}
	if !misc.IsListEqual(rems, rems2) {t.Errorf("At least one of the entries in the rems array is no longer equal after serialization")}

	// GetFromGob
	data, err := os.ReadFile(name)
	if err != nil {t.Errorf("Could not read from gob file: %v", err)}

	rems3, err := misc.GetFromGob[Entry](data)
	if err != nil {t.Errorf("Encountered an error while parsing the gob data: %v", err)}
	if !misc.IsListEqual(rems, rems3) {t.Errorf("At least one of the entries in the rems array is no longer equal after serialization and manual decoding")}
}
