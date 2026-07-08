package misc

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func Test_Filter(t *testing.T) {
	list1 := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	list2 := []int{}
	list3 := []int(nil)

	for _, res := range Filter(list1, func(entry int)bool {return entry % 2 == 0}) {
		if res % 2 != 0 {t.Errorf("Filter on list1 failed. Entry: %v, Expected res: %v%%2==0, Actual res: %v", res, res, res % 2)}
	}

	if l := len(Filter(list2, func(_ int)bool{return true})); l != 0 {
		t.Errorf("Got an unexpected length from filter on list2. Expected: 0, Got: %v", l)
	}

	if l := len(Filter(list3, func(_ int)bool{return true})); l != 0 {
		t.Errorf("Got an unexpected length from filter on list3. Expected: 0, Got: %v", l)
	}
}

type iwrap int
func (i iwrap) Compare(incoming iwrap) bool {return i == incoming}

func Test_IsListEqual(t *testing.T) {
	l1, l2, l3, l4 := []iwrap{1, 2, 3}, []iwrap{1, 2, 3}, []iwrap{1, 3, 2}, []iwrap{1, 2}

	if !IsListEqual(l1, l1) {t.Errorf("l1 and itself aren't equal when they should be")}
	if !IsListEqual(l1, l2) {t.Errorf("l1 and l2 aren't equal when they should be")}
	if  IsListEqual(l1, l3) {t.Errorf("l1 and l3 are equal when they shouldn't be")}
	if  IsListEqual(l1, l4) {t.Errorf("l1 and l4 are equal when they shouldn't be")}

	if !IsListEqual(l2, l2) {t.Errorf("l2 and itself aren't equal when they should be")}
	if !IsListEqual(l2, l1) {t.Errorf("l2 and l1 aren't equal when they should be")}
	if  IsListEqual(l2, l3) {t.Errorf("l2 and l3 are equal when they shouldn't be")}
	if  IsListEqual(l2, l4) {t.Errorf("l2 and l4 are equal when they shouldn't be")}

	if !IsListEqual(l3, l3) {t.Errorf("l3 and itself aren't equal when they should be")}
	if  IsListEqual(l3, l1) {t.Errorf("l3 and l1 are equal when they shouldn't be")}
	if  IsListEqual(l3, l2) {t.Errorf("l3 and l2 are equal when they shouldn't be")}
	if  IsListEqual(l3, l4) {t.Errorf("l3 and l4 are equal when they shouldn't be")}

	if !IsListEqual(l4, l4) {t.Errorf("l4 and itself aren't equal when they should be")}
	if  IsListEqual(l4, l1) {t.Errorf("l4 and l1 are equal when they shouldn't be")}
	if  IsListEqual(l4, l2) {t.Errorf("l4 and l2 are equal when they shouldn't be")}
	if  IsListEqual(l4, l3) {t.Errorf("l4 and l3 are equal when they shouldn't be")}

	// Checking every permutation is probably overkill for such a simple function, but hey it's something I can actually check
}

func Test_GetFileWatch(t *testing.T) {
	// Test that GetFileWatch errors correctly, and that it can get a file to watch
	dir := t.TempDir()
	dur, err := time.ParseDuration("5s")
	if err != nil {t.Errorf("Could not create duration object: %v", err)}

	NONEXISTENT := filepath.Join(dir, "NONEXISTENT")
	if _, err := GetFileWatch(NONEXISTENT, dur); err == nil {
		t.Errorf("Was able to create a file watch for a file which should not exist")
	}

	// Not sure how to write a test for the stat emitting an error
}

func Test_WatchedFile_Close(t *testing.T) {
	// Check that WatchedFile.Close doesn't break anything like an active write

	var wf *WatchedFile = nil
	if err := wf.Close(); err == nil {
		t.Errorf("Didn't get an error when trying to close a nil WatchedFile pointer")
	}

	if err := (&WatchedFile{}).Close(); err == nil {
		t.Errorf("Didn't get an error from closing an empty WatchedFile")
	}

	if err := (&WatchedFile{closed: true}).Close(); err != nil {
		t.Errorf("Got an error from a WatchedFile that should be already closed: %v", err)
	}

	// Again, don't know how to test an interrupting close
}

func Test_WatchedFile_CheckForUpdate(t *testing.T) {
	// Check that a write gets properly noticed

	var wf *WatchedFile = nil
	if _, err := wf.CheckForUpdate(); err == nil {
		t.Errorf("Didn't get an error when trying to check for an update on a nil WatchedFile pointer")
	}

	if _, err := (&WatchedFile{closed: true}).CheckForUpdate(); err == nil {
		t.Errorf("Was able to check for an update on an already closed WatchedFile")
	}

	// Setup
	dir := t.TempDir()

	dur1, err := time.ParseDuration("1s")
	if err != nil {t.Errorf("Could not create duration object: %v", err)}
	dur2, err := time.ParseDuration("2.5s")
	if err != nil {t.Errorf("Could not create duration object: %v", err)}


	name := filepath.Join(dir, "watch.txt")
	file, err := os.Create(name)
	if err != nil {t.Errorf("Could not open \"%v\" for writing: %v", name, err)}

	watch, err := GetFileWatch(name, dur1)
	if err != nil {t.Errorf("Could not create watch for \"%v\": %v", name, err)}

	updated := make(chan int)
	// End setup



	// Write to file occasionally
	go func(){
		for i := range 5 {
			file.WriteString("Writing iteration " + fmt.Sprint(i) + "\n")
			time.Sleep(dur2)
		}
	}()

	// Watch the file for 3 updates
	go func(){
		u := 0
		for i := 0; i < 30 && u < 3; i++ { // stop after 30 seconds or 3 updates
			res, err := watch.CheckForUpdate()
			if err != nil {t.Errorf("Got an error while checking for an update: %v", err)}
			if res {u++}
			watch.Sleep()
		}

		updated <- u
	}()

	if i := <-updated; i != 3 {
		t.Errorf("Got an unexpected number of file updates. Expected: 3, Got: %v", i)
	}


	if err := os.RemoveAll(dir); err != nil {
		fmt.Printf("Warning: Could not remove temp dir \"%v\": %v", dir, err)
	}
}
