package main

import (
	"os"
	"testing"
)

func Test__ReadFileResult_new(t *testing.T) {
	t1 := *ReadFileResult_new("", nil, nil, nil, "");
	if t1.keyword != "" {t.Errorf("t1: Mismatched keyword")}
	if t1.emsg != nil {t.Errorf("t1: Mismatched error")}
	if t1.result != "" {t.Errorf("t1: Mismatched result")}

	t1res := t1.callback(make([]string, 0))
	if !t1res.Deep_compare(t1) {t.Errorf("t1: Deep compare returned false")}


	inputs := []string{"testing", "value1", "value2"}

	// A tad bit hacky, but if it works it works
	var callback func(inputs []string) ReadFileResult;
	callback = func(inputs []string) ReadFileResult {
		return *ReadFileResult_new(
			inputs[0],
			inputs,
			callback,
			nil,
			inputs[len(inputs)-1],
		)
	}

	t2 := *ReadFileResult_new("testing", inputs, callback, nil, "value2");
	if t2.keyword != "testing" {t.Errorf("t2: Mismatched keyword")}
	if t2.emsg != nil {t.Errorf("t2: Mismatched error")}
	if t2.result != "value2" {t.Errorf("t2: Mismatched result")}

	t2res := t2.callback(inputs)
	if !t2res.Deep_compare(t2) {t.Errorf("t2: Deep compare returned false")}
}

func Test__ReadFile(t *testing.T) {

}

func writeConfigFile() (*os.File, error) {
	file, err := os.Create("config.trem.test")
	if err != nil {return nil, err}

	// write file contents

	_, err = file.Seek(0, 0)
	if err != nil {return nil, err}
	return file, nil
}
func Test__ReadConfigFile(t *testing.T) {
	file, confErr := writeConfigFile()
	if confErr != nil {t.Errorf("%v", confErr)}
	defer file.Close()

	var funcmap map[string]func([]string) ReadFileResult = make(map[string]func([]string) ReadFileResult) // fill out later
	results, readErr := ReadFile(file, funcmap)
	if readErr != nil {t.Errorf("%v", readErr)}

	// Check the results list for correctness
	for idx := 0; idx < len(results); idx++ {

	}
}

func writeReminderFile() (*os.File, error) {
	file, err := os.Create("reminder.trem.test")
	if err != nil {return nil, err}

	// write file contents

	_, err = file.Seek(0, 0)
	if err != nil {return nil, err}
	return file, nil
}
func Test__ReadReminderFile(t *testing.T) {
	file, remErr := writeReminderFile()
	if remErr != nil {t.Errorf("%v", remErr)}
	defer file.Close()
}