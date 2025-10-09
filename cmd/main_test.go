package main

import (
	"encoding/json"
	"testing"
)

func TestValidateFlags(t *testing.T) {
	tests := []struct {
		name      string
		eventType string
		message   string
		wantErr   bool
	}{
		{
			name:      "valid task type with message",
			eventType: "task",
			message:   "Test message",
			wantErr:   false,
		},
		{
			name:      "valid question type with message",
			eventType: "question",
			message:   "Test question",
			wantErr:   false,
		},
		{
			name:      "valid progress type with message",
			eventType: "progress",
			message:   "50% complete",
			wantErr:   false,
		},
		{
			name:      "valid session type with message",
			eventType: "session",
			message:   "Session started",
			wantErr:   false,
		},
		{
			name:      "missing event type",
			eventType: "",
			message:   "Test message",
			wantErr:   true,
		},
		{
			name:      "missing message",
			eventType: "task",
			message:   "",
			wantErr:   true,
		},
		{
			name:      "invalid event type",
			eventType: "invalid",
			message:   "Test message",
			wantErr:   true,
		},
		{
			name:      "both missing",
			eventType: "",
			message:   "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFlags(tt.eventType, tt.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateFlags() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMapSubject(t *testing.T) {
	tests := []struct {
		name      string
		eventType string
		state     string
		want      string
	}{
		// Task mappings
		{
			name:      "task in_progress",
			eventType: "task",
			state:     "in_progress",
			want:      "claude.tasks.started",
		},
		{
			name:      "task completed",
			eventType: "task",
			state:     "completed",
			want:      "claude.tasks.completed",
		},
		{
			name:      "task blocked",
			eventType: "task",
			state:     "blocked",
			want:      "claude.tasks.blocked",
		},
		{
			name:      "task with pending state",
			eventType: "task",
			state:     "pending",
			want:      "claude.tasks",
		},
		{
			name:      "task without state",
			eventType: "task",
			state:     "",
			want:      "claude.tasks",
		},
		// Question mappings
		{
			name:      "question blocked",
			eventType: "question",
			state:     "blocked",
			want:      "claude.questions.waiting",
		},
		{
			name:      "question without state",
			eventType: "question",
			state:     "",
			want:      "claude.questions.asked",
		},
		{
			name:      "question with any other state",
			eventType: "question",
			state:     "answered",
			want:      "claude.questions.asked",
		},
		// Progress mappings
		{
			name:      "progress",
			eventType: "progress",
			state:     "",
			want:      "claude.progress.update",
		},
		{
			name:      "progress with state",
			eventType: "progress",
			state:     "50",
			want:      "claude.progress.update",
		},
		// Session mappings
		{
			name:      "session completed",
			eventType: "session",
			state:     "completed",
			want:      "claude.session.completed",
		},
		{
			name:      "session started",
			eventType: "session",
			state:     "",
			want:      "claude.session.started",
		},
		{
			name:      "session with any other state",
			eventType: "session",
			state:     "active",
			want:      "claude.session.started",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapSubject(tt.eventType, tt.state)
			if got != tt.want {
				t.Errorf("mapSubject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessageJSONMarshaling(t *testing.T) {
	tests := []struct {
		name    string
		msg     Message
		wantErr bool
	}{
		{
			name: "complete message",
			msg: Message{
				Event:     "claude.tasks.started",
				Timestamp: "2025-10-09T14:30:00Z",
				SessionID: "test-session",
				Message:   "Test message",
				State:     "in_progress",
				TaskNum:   "3/15",
			},
			wantErr: false,
		},
		{
			name: "message with empty optional fields",
			msg: Message{
				Event:     "claude.progress.update",
				Timestamp: "2025-10-09T14:30:00Z",
				Message:   "50% complete",
			},
			wantErr: false,
		},
		{
			name: "message with only required fields",
			msg: Message{
				Event:     "claude.questions.asked",
				Timestamp: "2025-10-09T14:30:00Z",
				Message:   "What should I do?",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("json.Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				// Verify we can unmarshal it back
				var unmarshaled Message
				if err := json.Unmarshal(data, &unmarshaled); err != nil {
					t.Errorf("json.Unmarshal() error = %v", err)
				}

				// Verify required fields match
				if unmarshaled.Event != tt.msg.Event {
					t.Errorf("Event mismatch: got %v, want %v", unmarshaled.Event, tt.msg.Event)
				}
				if unmarshaled.Message != tt.msg.Message {
					t.Errorf("Message mismatch: got %v, want %v", unmarshaled.Message, tt.msg.Message)
				}
			}
		})
	}
}

func TestMessageOmitEmpty(t *testing.T) {
	msg := Message{
		Event:     "claude.tasks.started",
		Timestamp: "2025-10-09T14:30:00Z",
		Message:   "Test message",
		// SessionID, State, and TaskNum are empty
	}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	// Parse the JSON to check that empty fields are omitted
	var rawJSON map[string]any
	if err := json.Unmarshal(data, &rawJSON); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	// Check that optional fields are not present in JSON
	if _, exists := rawJSON["session_id"]; exists {
		t.Error("session_id should be omitted when empty")
	}
	if _, exists := rawJSON["state"]; exists {
		t.Error("state should be omitted when empty")
	}
	if _, exists := rawJSON["task_num"]; exists {
		t.Error("task_num should be omitted when empty")
	}

	// Check that required fields are present
	if _, exists := rawJSON["event"]; !exists {
		t.Error("event should be present")
	}
	if _, exists := rawJSON["timestamp"]; !exists {
		t.Error("timestamp should be present")
	}
	if _, exists := rawJSON["message"]; !exists {
		t.Error("message should be present")
	}
}

func TestPrintSuccess(t *testing.T) {
	// This is a simple test to ensure printSuccess doesn't panic
	// In a real-world scenario, you might want to capture stdout
	tests := []struct {
		name      string
		eventType string
		state     string
		taskNum   string
		sessionID string
		subject   string
	}{
		{
			name:      "full output",
			eventType: "task",
			state:     "in_progress",
			taskNum:   "3/15",
			sessionID: "test-session",
			subject:   "claude.tasks.started",
		},
		{
			name:      "minimal output",
			eventType: "progress",
			state:     "",
			taskNum:   "",
			sessionID: "",
			subject:   "claude.progress.update",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just ensure it doesn't panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("printSuccess() panicked: %v", r)
				}
			}()
			printSuccess(tt.eventType, tt.state, tt.taskNum, tt.sessionID, tt.subject)
		})
	}
}

func TestExitCodes(t *testing.T) {
	// Verify exit code constants are defined correctly
	if exitSuccess != 0 {
		t.Errorf("exitSuccess should be 0, got %d", exitSuccess)
	}
	if exitInvalidArgs != 1 {
		t.Errorf("exitInvalidArgs should be 1, got %d", exitInvalidArgs)
	}
	if exitConnectionError != 2 {
		t.Errorf("exitConnectionError should be 2, got %d", exitConnectionError)
	}
}

func TestValidTypes(t *testing.T) {
	expectedTypes := []string{"task", "question", "progress", "session"}

	for _, typ := range expectedTypes {
		if !validTypes[typ] {
			t.Errorf("validTypes should contain '%s'", typ)
		}
	}

	// Check that unexpected types are not valid
	invalidTypes := []string{"invalid", "foo", "bar", ""}
	for _, typ := range invalidTypes {
		if validTypes[typ] {
			t.Errorf("validTypes should not contain '%s'", typ)
		}
	}
}

func TestAuthTypeDefaults(t *testing.T) {
	// Test that default auth type is set correctly
	if defaultAuthType != "none" {
		t.Errorf("defaultAuthType should be 'none', got '%s'", defaultAuthType)
	}

	// Test that default URL is set correctly
	expectedURL := "nats://localhost:4222"
	if defaultNATSURL != expectedURL {
		t.Errorf("defaultNATSURL should be '%s', got '%s'", expectedURL, defaultNATSURL)
	}

	// Test that credential fields are empty by default
	if defaultUsername != "" {
		t.Error("defaultUsername should be empty")
	}
	if defaultPassword != "" {
		t.Error("defaultPassword should be empty")
	}
	if defaultToken != "" {
		t.Error("defaultToken should be empty")
	}
	if defaultNKey != "" {
		t.Error("defaultNKey should be empty")
	}
	if defaultNATSJWT != "" {
		t.Error("defaultNATSJWT should be empty")
	}
	if defaultNATSSeed != "" {
		t.Error("defaultNATSSeed should be empty")
	}
}

func TestConnectNATSWithEnvVars(t *testing.T) {
	// Save original env vars
	origURL := defaultNATSURL
	origAuthType := defaultAuthType
	defer func() {
		defaultNATSURL = origURL
		defaultAuthType = origAuthType
	}()

	tests := []struct {
		name        string
		authType    string
		username    string
		password    string
		token       string
		nkey        string
		jwt         string
		seed        string
		expectError bool
	}{
		{
			name:        "none authentication",
			authType:    "none",
			expectError: true, // Will fail without NATS server
		},
		{
			name:        "userpass authentication",
			authType:    "userpass",
			username:    "testuser",
			password:    "testpass",
			expectError: true, // Will fail without NATS server
		},
		{
			name:        "token authentication",
			authType:    "token",
			token:       "testtoken",
			expectError: true, // Will fail without NATS server
		},
		{
			name:        "nkey authentication",
			authType:    "nkey",
			nkey:        "SUAKTEST",
			expectError: true, // Will fail without NATS server
		},
		{
			name:        "decentralized authentication",
			authType:    "decentralized",
			jwt:         "test.jwt.token",
			seed:        "SUAKTEST",
			expectError: true, // Will fail without NATS server
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up test configuration
			defaultAuthType = tt.authType
			defaultUsername = tt.username
			defaultPassword = tt.password
			defaultToken = tt.token
			defaultNKey = tt.nkey
			defaultNATSJWT = tt.jwt
			defaultNATSSeed = tt.seed
			defaultNATSURL = "nats://localhost:14222" // Use non-standard port

			// This will fail because there's no NATS server, but we're testing
			// that the function constructs the connection attempt correctly
			_, err := connectNATS()

			// We expect an error because no NATS server is running
			if err == nil && tt.expectError {
				t.Error("Expected error due to no NATS server, but got nil")
			}
		})
	}
}
