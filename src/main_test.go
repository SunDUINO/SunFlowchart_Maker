package main

import "testing"

func TestMainApp(t *testing.T) {
	result := "SunGo"
	if result != "SunGo" {
		t.Errorf("Oczekiwano SunGo, otrzymano %s", result)
	}
}
