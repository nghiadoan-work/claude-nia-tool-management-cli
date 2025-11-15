package ui

import (
	"os"
	"time"

	"github.com/briandowns/spinner"
)

// Spinner wraps the briandowns/spinner with our custom settings
type Spinner struct {
	s *spinner.Spinner
}

// NewSpinner creates a new spinner with a message
func NewSpinner(message string) *Spinner {
	s := spinner.New(
		spinner.CharSets[14], // Use the dots spinner
		100*time.Millisecond,
		spinner.WithWriter(os.Stderr),
	)
	s.Suffix = " " + message
	s.Color("cyan", "bold")

	return &Spinner{s: s}
}

// Start starts the spinner
func (sp *Spinner) Start() {
	sp.s.Start()
}

// Stop stops the spinner
func (sp *Spinner) Stop() {
	sp.s.Stop()
}

// UpdateMessage updates the spinner message
func (sp *Spinner) UpdateMessage(message string) {
	sp.s.Suffix = " " + message
}

// Success stops the spinner and shows a success message
func (sp *Spinner) Success(message string) {
	sp.s.Stop()
	PrintSuccess("%s", message)
}

// Fail stops the spinner and shows an error message
func (sp *Spinner) Fail(message string) {
	sp.s.Stop()
	PrintError("%s", message)
}

// WithSpinner executes a function while showing a spinner
// Returns the spinner for further manipulation if needed
func WithSpinner(message string, fn func() error) error {
	sp := NewSpinner(message)
	sp.Start()
	defer sp.Stop()

	err := fn()
	if err != nil {
		sp.Fail(message + " failed")
		return err
	}

	sp.Success(message + " completed")
	return nil
}

// SpinnerFunc is a function that can be executed with a spinner
type SpinnerFunc func() error

// Execute runs the function with a spinner showing the given message
func (sf SpinnerFunc) Execute(message string) error {
	return WithSpinner(message, sf)
}
