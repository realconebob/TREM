package main

import (
	"testing"
)

func Test__Trem_ReadFileResult_new(t *testing.T) {
	t1 := *Trem_ReadFileResult_new("", nil, nil, nil, "");
	if t1.keyword != "" {t.Errorf("t1: Mismatched keyword")}
	if t1.emsg != nil {t.Errorf("t1: Mismatched error")}
	if t1.result != "" {t.Errorf("t1: Mismatched result")}

	t1res := t1.callback(make([]string, 0))
	if !t1res.Deep_compare(t1) {t.Errorf("t1: Deep compare returned false")}


	inputs := []string{"testing", "value1", "value2"}

	// A tad bit hacky, but if it works it works
	var callback func(inputs []string) Trem_ReadFileResult;
	callback = func(inputs []string) Trem_ReadFileResult {
		return *Trem_ReadFileResult_new(
			inputs[0],
			inputs,
			callback,
			nil,
			inputs[len(inputs)-1],
		)
	}

	t2 := *Trem_ReadFileResult_new("testing", inputs, callback, nil, "value2");
	if t2.keyword != "testing" {t.Errorf("t2: Mismatched keyword")}
	if t2.emsg != nil {t.Errorf("t2: Mismatched error")}
	if t2.result != "value2" {t.Errorf("t2: Mismatched result")}

	t2res := t2.callback(inputs)
	if !t2res.Deep_compare(t2) {t.Errorf("t2: Deep compare returned false")}
}

func Test__Trem_ReadFile(t *testing.T) {

}

// TODO: Test ActOnFile