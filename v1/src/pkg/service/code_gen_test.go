package service

import "testing"

func TestMakeURL(t *testing.T) {
	ur := MakeURL("")
	if ur == "" {
		t.Fail()
	}
}
