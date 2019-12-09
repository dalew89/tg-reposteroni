package main

import "testing"

// Tests for IdentifyMessage
func TestIncomingMessage_IdentifyMessage(t *testing.T) {
	// Test for a message with URL
	messageWithURL := IncomingMessage{MessageText: "This string has a URL https://go.dev"}
	result := messageWithURL.IdentifyMessage()
	if result == "" {
		t.Errorf("Expected a url, got a blank field")
	}

	// Test for a message without URL
	messageWithoutURL := IncomingMessage{MessageText: "This string without URL"}
	value := messageWithoutURL.IdentifyMessage()
	if value != "" {
		t.Errorf("Expected an empty string, got %s", value)
	}
}

