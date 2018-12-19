package service

import "testing"

func TestMakeUrl(t *testing.T) {
	ur := MakeUrl("oumamarvin@gmail.com")
	if ur == "" {
		t.Fail(	)
	}
}
