package config

import (
	"testing"
)

func TestConfig_AddCloser(t *testing.T) {
	config := &Config{}

	// Create a mock closer function
	mockCloser := func() {}

	// Call the AddCloser method
	config.AddCloser(mockCloser)

	// Check if the closer was added correctly
	if len(config.closers) != 1 {
		t.Errorf("AddCloser failed to add the closer function")
	}
}

func TestConfig_Close(t *testing.T) {
	config := &Config{}

	testValue := 0
	// Create a mock closer function
	mockCloser := func() {
		testValue = 1
	}

	// Add the mock closer to the Config
	config.AddCloser(mockCloser)

	// Call the Close method
	config.Close()

	// Check if the closer was called
	if testValue != 1 {
		t.Errorf("Close failed to call the closer function")
	}
}
