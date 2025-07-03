package main

import (
	"fmt"
	"os"
)

type EchoWaveError struct {
	Operation string
	Err       error
}

func (e *EchoWaveError) Error() string {
	return fmt.Sprintf("%s failed: %v", e.Operation, e.Err)
}

func newError(operation string, err error) *EchoWaveError {
	return &EchoWaveError{
		Operation: operation,
		Err:       err,
	}
}

func exitWithError(err error) {
	fmt.Printf("❌ %v\n", err)
	os.Exit(1)
}

func exitWithMessage(message string) {
	fmt.Printf("❌ %s\n", message)
	os.Exit(1)
}