package ui

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewSpinner(t *testing.T) {
	sp := NewSpinner("test message")
	assert.NotNil(t, sp)
	assert.NotNil(t, sp.s)
}

func TestSpinnerStartStop(t *testing.T) {
	sp := NewSpinner("test")
	assert.NotPanics(t, func() {
		sp.Start()
		time.Sleep(100 * time.Millisecond)
		sp.Stop()
	})
}

func TestSpinnerUpdateMessage(t *testing.T) {
	sp := NewSpinner("initial message")
	assert.NotPanics(t, func() {
		sp.UpdateMessage("updated message")
	})
}

func TestSpinnerSuccess(t *testing.T) {
	sp := NewSpinner("test")
	sp.Start()
	assert.NotPanics(t, func() {
		sp.Success("operation successful")
	})
}

func TestSpinnerFail(t *testing.T) {
	sp := NewSpinner("test")
	sp.Start()
	assert.NotPanics(t, func() {
		sp.Fail("operation failed")
	})
}

func TestWithSpinner_Success(t *testing.T) {
	callCount := 0
	err := WithSpinner("test operation", func() error {
		callCount++
		time.Sleep(50 * time.Millisecond)
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, callCount)
}

func TestWithSpinner_Error(t *testing.T) {
	expectedErr := errors.New("test error")
	err := WithSpinner("test operation", func() error {
		return expectedErr
	})

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestSpinnerFunc_Execute(t *testing.T) {
	var executed bool
	fn := SpinnerFunc(func() error {
		executed = true
		return nil
	})

	err := fn.Execute("test message")
	assert.NoError(t, err)
	assert.True(t, executed)
}

func TestSpinnerFunc_ExecuteWithError(t *testing.T) {
	expectedErr := errors.New("execution error")
	fn := SpinnerFunc(func() error {
		return expectedErr
	})

	err := fn.Execute("test message")
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestSpinnerConcurrency(t *testing.T) {
	// Test that multiple spinners don't interfere
	sp1 := NewSpinner("spinner 1")
	sp2 := NewSpinner("spinner 2")

	assert.NotPanics(t, func() {
		sp1.Start()
		sp2.Start()
		time.Sleep(50 * time.Millisecond)
		sp1.Stop()
		sp2.Stop()
	})
}
