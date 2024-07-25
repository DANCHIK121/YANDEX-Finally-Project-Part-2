package main

import (
	"testing"
)

func TestReadJson(t *testing.T) {
	testCases := TestStructForReadJson{
		out: JsonConstsFile{
			// Check current json consts file
			ServerIP: "127.0.0.1",
			ServerPort: "8000",
			ServerJWTTokenEXP: 20,
			MaxBytesForReader: 1048576,
		},
    }

	result := ReadJson()

	if result != testCases.out {
		t.Error()
	}
}