package main

import (
	"testing"
)

func TestCalculatingResult(t *testing.T) {
	EnvironmentVariablesInit()

	testCases := []TestStructForCalculatingResult{
		{in: Task{ID: 0, Arg1: 5, Arg2: 2, Operation: "-"}, out: 3},   // 1
		{in: Task{ID: 1, Arg1: 5, Arg2: 2, Operation: "+"}, out: 7},   // 2
		{in: Task{ID: 2, Arg1: 10, Arg2: 2, Operation: "*"}, out: 20}, // 3
		{in: Task{ID: 2, Arg1: 10, Arg2: 2, Operation: "/"}, out: 5},  // 4
    }

	for test := range testCases {
        result := CalculatingResult(testCases[test].in)

        if result != testCases[test].out {
            t.Error()
        }
    }
}