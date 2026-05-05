package main

import (
	"testing"
)

func Test_aofResult(t *testing.T) {
	val := *New_AOFResult[string]()

	if val.line != "" || val.oper == nil || val.res != "" {
		t.Errorf("val was not \"emptied\" correctly")
	}

	if val.oper(val.line) != "" {
		t.Error("val's empty oper does not return the correct value")
	}
}

// TODO: Test ActOnFile