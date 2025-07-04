package main

import (
	"fmt"
	"os"
)

// EchoWaveError wraps errors with operation context for better debugging.
// Provides consistent error formatting throughout the EchoWave application.
type EchoWaveError struct {
	Operation string
	Err       error
}

// Error implements the error interface for EchoWaveError by returning a formatted
// error message that combines the operation name and underlying error details.
// This method provides consistent error formatting across the application with
// the pattern "[operation] failed: [error details]".
func (e *EchoWaveError) Error() string {
	return fmt.Sprintf("%s failed: %v", e.Operation, e.Err)
}

// newError creates a new EchoWaveError instance that wraps an underlying error
// with contextual information about the operation that failed. The operation parameter
// should describe what was being attempted when the error occurred, while err contains
// the original error details. This function provides consistent error wrapping throughout
// the application for better debugging and user feedback.
func newError(operation string, err error) *EchoWaveError {
	return &EchoWaveError{
		Operation: operation,
		Err:       err,
	}
}

// exitWithError prints a formatted error message to stdout with an error emoji
// and immediately terminates the program with exit code 1. This function provides
// a consistent way to handle fatal errors throughout the application by displaying
// user-friendly error messages before graceful shutdown.
func exitWithError(err error) {
	fmt.Printf("❌ %v\n", err)
	os.Exit(1)
}
