package main

import (
	"testing"
)

func TestReadJson(t *testing.T) {
	testCases := TestStructForReadJson{
		out: JsonConstsFile{
			// Check current json consts file
			ComputingPower: 4,
			TimeAdditionMS: 1000,
			TimeSubstrationMS: 1000,
			TimeMultiplicationsMS: 1000,
			TimeDivisionsMS: 1000,
		},
    }

	result := ReadJson()

	if result != testCases.out {
		t.Error()
	}
}
