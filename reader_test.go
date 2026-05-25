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

func writeConfigFile() (*os.File, error) {
	file, err := os.Create("config.trem.test")
	if err != nil {return nil, err}

	file.Write([]byte(
		"type=config\n" +
		"dirs=%HOME/.text-reminders;\"/di re ct or y\"\n",
	))

	_, err = file.Seek(0, 0)
	if err != nil {
		file.Close()
		return nil, err
	}
	return file, nil
}
func Test__ReadConfigFile(t *testing.T) {
	file, confErr := writeConfigFile()
	if confErr != nil {t.Errorf("%v", confErr)}
	defer file.Close()

	var funcmap map[string]func([]string) ReadFileResult = make(map[string]func([]string) ReadFileResult) // fill out later
	results, readErr := ReadFile(file, funcmap)
	if readErr != nil {t.Errorf("%v", readErr)}
	if len(results) != 2 {t.Errorf("Results mismatch. Expected 2 results, got: %v", len(results))}

	// Check the results list for correctness
	for idx, res := range results {
		switch idx {
		case 0:
			if res.inputs[0] != "type" {t.Errorf("keyword mismatch. Expected: type. Got: " + res.inputs[0])}
			if res.inputs[1] != "config" {t.Errorf("input mismatch. Expected: config. Got: " + res.inputs[1])}

		case 1:
			if res.inputs[0] != "dirs" {t.Errorf("keyword mismatch. Expected: . Got: " + res.inputs[0])}
			if res.inputs[1] != "%HOME/.text-reminders;\"/di re ct or y\"" {t.Errorf("input mismatch. Expected: " + "%HOME/.text-reminders;\"/di re ct or y\"" + ". Got: " + res.inputs[1])}

		default:
			t.Errorf("Broke out of switch case")
		}
	}
}

func writeReminderFile() (*os.File, error) {
	file, err := os.Create("reminder.trem.test")
	if err != nil {return nil, err}

	file.Write([]byte(
		"type=reminder\n" +
		"date=monday\n" +
		"reminder=Make sure to check your email\n",
	))

	_, err = file.Seek(0, 0)
	if err != nil {
		file.Close()
		return nil, err
	}
	return file, nil
}
func Test__ReadReminderFile(t *testing.T) {
	file, remErr := writeReminderFile()
	if remErr != nil {t.Errorf("%v", remErr)}
	defer file.Close()

	var funcmap map[string]func([]string) ReadFileResult = make(map[string]func([]string) ReadFileResult) // fill out later
	results, readErr := ReadFile(file, funcmap)
	if readErr != nil {t.Errorf("%v", readErr)}
	if len(results) != 3 {t.Errorf("Result mismatch. Expected 3 results, got: %v", len(results))}

	for idx, res := range results {
		switch idx {
		case 0:
			if res.inputs[0] != "" {t.Errorf("keyword mismatch. Expected: type. Got: " + res.inputs[0])}
			if res.inputs[1] != "" {t.Errorf("input mismatch. Expected: reminder. Got: " + res.inputs[1])}

		case 1:
			if res.inputs[0] != "" {t.Errorf("keyword mismatch. Expected: date. Got: " + res.inputs[0])}
			if res.inputs[1] != "" {t.Errorf("input mismatch. Expected: monday. Got: " + res.inputs[1])}

		case 2:
			if res.inputs[0] != "" {t.Errorf("keyword mismatch. Expected: reminder. Got: " + res.inputs[0])}
			if res.inputs[1] != "" {t.Errorf("input mismatch. Expected: Make sure to check your email. Got: " + res.inputs[1])}

		default:
			t.Errorf("")
		}
	}
}