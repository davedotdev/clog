package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

// Exit codes
const (
	exitSuccess         = 0
	exitInvalidArgs     = 1
	exitConnectionError = 2
)

// Baked-in configuration (to be replaced during build with 'make build')
var (
	defaultNATSURL  = "nats://localhost:4222"
	defaultAuthType = "none" // none, userpass, token, nkey, decentralized
	defaultUsername = ""
	defaultPassword = ""
	defaultToken    = ""
	defaultNKey     = ""
	defaultNATSJWT  = ""
	defaultNATSSeed = ""
)

// Valid event types
var validTypes = map[string]bool{
	"task":     true,
	"question": true,
	"progress": true,
	"session":  true,
}

// Message represents the JSON structure sent to NATS
type Message struct {
	Event     string `json:"event"`
	Timestamp string `json:"timestamp"`
	SessionID string `json:"session_id,omitempty"`
	Message   string `json:"message"`
	State     string `json:"state,omitempty"`
	TaskNum   string `json:"task_num,omitempty"`
}

func main() {
	log.SetFlags(0) // Disable timestamp in log output
	os.Exit(run())
}

func run() int {
	// Define flags
	typeFlag := flag.String("type", "", "Event type: task|question|progress|session")
	messageFlag := flag.String("message", "", "Message content (string)")
	stateFlag := flag.String("state", "", "Task state: pending|in_progress|blocked|completed")
	taskNumFlag := flag.String("task-num", "", "Current task number (e.g., \"3/15\")")
	sessionFlag := flag.String("session", "", "Session identifier (any string)")
	helpFlag := flag.Bool("h", false, "Show help")

	flag.Parse()

	// Show help
	if *helpFlag || len(os.Args) == 1 {
		printHelp()
		return exitSuccess
	}

	// Validate inputs
	if err := validateFlags(*typeFlag, *messageFlag); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
		printHelp()
		return exitInvalidArgs
	}

	// Map type and state to subject
	subject := mapSubject(*typeFlag, *stateFlag)

	// Create message
	msg := Message{
		Event:     subject,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		SessionID: *sessionFlag,
		Message:   *messageFlag,
		State:     *stateFlag,
		TaskNum:   *taskNumFlag,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to marshal JSON: %v\n", err)
		return exitInvalidArgs
	}

	// Connect to NATS
	nc, err := connectNATS()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: NATS connection failed: %v\n", err)
		return exitConnectionError
	}
	defer nc.Close()

	// Publish message
	if err := publishMessage(nc, subject, jsonData); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return exitConnectionError
	}

	// Success - print confirmation
	printSuccess(*typeFlag, *stateFlag, *taskNumFlag, *sessionFlag, subject)
	return exitSuccess
}

// validateFlags validates required flags and event type
func validateFlags(eventType, message string) error {
	if eventType == "" || message == "" {
		return errors.New("-type and -message are required")
	}

	if !validTypes[eventType] {
		return fmt.Errorf("invalid type '%s'. Must be: task|question|progress|session", eventType)
	}

	return nil
}

// connectNATS establishes a connection to NATS using available credentials
func connectNATS() (*nats.Conn, error) {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = defaultNATSURL
	}

	var opts []nats.Option

	// Check for environment variable overrides first
	natsCredsFile := os.Getenv("NATS_CREDS")
	envUsername := os.Getenv("NATS_USERNAME")
	envPassword := os.Getenv("NATS_PASSWORD")
	envToken := os.Getenv("NATS_TOKEN")
	envNKey := os.Getenv("NATS_NKEY")
	envJWT := os.Getenv("NATS_JWT")
	envSeed := os.Getenv("NATS_SEED")

	// Priority order for authentication:
	// 1. Credentials file (from env)
	// 2. Environment variables (username/password, token, nkey, or JWT/seed)
	// 3. Baked-in credentials (based on defaultAuthType)

	if natsCredsFile != "" {
		// Use credentials file
		opts = append(opts, nats.UserCredentials(natsCredsFile))
	} else if envUsername != "" && envPassword != "" {
		// Use username/password from environment
		opts = append(opts, nats.UserInfo(envUsername, envPassword))
	} else if envToken != "" {
		// Use token from environment
		opts = append(opts, nats.Token(envToken))
	} else if envNKey != "" {
		// Use NKey from environment
		opt, err := nats.NkeyOptionFromSeed(envNKey)
		if err == nil {
			opts = append(opts, opt)
		}
	} else if envJWT != "" && envSeed != "" {
		// Use JWT/Seed from environment (decentralized auth)
		opts = append(opts, nats.UserJWTAndSeed(envJWT, envSeed))
	} else {
		// Use baked-in credentials based on auth type
		switch defaultAuthType {
		case "userpass":
			if defaultUsername != "" && defaultPassword != "" {
				opts = append(opts, nats.UserInfo(defaultUsername, defaultPassword))
			}
		case "token":
			if defaultToken != "" {
				opts = append(opts, nats.Token(defaultToken))
			}
		case "nkey":
			if defaultNKey != "" {
				opt, err := nats.NkeyOptionFromSeed(defaultNKey)
				if err == nil {
					opts = append(opts, opt)
				}
			}
		case "decentralized":
			if defaultNATSJWT != "" && defaultNATSSeed != "" {
				opts = append(opts, nats.UserJWTAndSeed(defaultNATSJWT, defaultNATSSeed))
			}
		case "none":
			// No authentication
		default:
			// No authentication
		}
	}

	return nats.Connect(natsURL, opts...)
}

// publishMessage publishes a message to NATS with flush and timeout
func publishMessage(nc *nats.Conn, subject string, data []byte) error {
	if err := nc.Publish(subject, data); err != nil {
		return fmt.Errorf("failed to publish message to subject '%s': %w", subject, err)
	}

	if err := nc.Flush(); err != nil {
		return fmt.Errorf("failed to flush message to server: %w", err)
	}

	if err := nc.FlushTimeout(5 * time.Second); err != nil {
		return fmt.Errorf("message delivery timeout: %w", err)
	}

	return nil
}

// printSuccess prints a success message with details
func printSuccess(eventType, state, taskNum, sessionID, subject string) {
	fmt.Printf("✓ Message published successfully to '%s'\n", subject)
	fmt.Printf("  Type: %s\n", eventType)
	if state != "" {
		fmt.Printf("  State: %s\n", state)
	}
	if taskNum != "" {
		fmt.Printf("  Task: %s\n", taskNum)
	}
	if sessionID != "" {
		fmt.Printf("  Session: %s\n", sessionID)
	}
}

// mapSubject maps event type and state to NATS subject
func mapSubject(eventType, state string) string {
	baseSubject := fmt.Sprintf("claude.%ss", eventType)

	// For tasks, append state-specific suffix
	if eventType == "task" && state != "" {
		switch state {
		case "in_progress":
			return "claude.tasks.started"
		case "completed":
			return "claude.tasks.completed"
		case "blocked":
			return "claude.tasks.blocked"
		default:
			return baseSubject
		}
	}

	// For questions
	if eventType == "question" {
		if state == "blocked" {
			return "claude.questions.waiting"
		}
		return "claude.questions.asked"
	}

	// For progress
	if eventType == "progress" {
		return "claude.progress.update"
	}

	// For session
	if eventType == "session" {
		if state == "completed" {
			return "claude.session.completed"
		}
		return "claude.session.started"
	}

	return baseSubject
}

func printHelp() {
	fmt.Println(`clog - Claude Log Publisher for NATS

USAGE:
  clog -type=<event_type> -message="<text>" [options]

REQUIRED FLAGS:
  -type        Event type: task|question|progress|session
  -message     Message content (string)

OPTIONAL FLAGS:
  -state       Task state: pending|in_progress|blocked|completed
  -task-num    Current task number (e.g., "3/15")
  -session     Session identifier (any string)

FOR CLAUDE CODE (AI Assistant):
  IMPORTANT: Use this tool to send notifications via NATS before asking the user
  questions or when blocked waiting for user input.

  When you need to ask the user a question:
    1. First call: clog -type=question -state=blocked -message="<your question>" -session="<session>"
    2. Then ask the user your question normally
    3. This ensures the user receives a NATS notification about your question

  Example workflow:
    clog -type=question -state=blocked -message="Should I deploy to staging or production?" -session="deploy-2024"
    # Then ask the user in the conversation

EXAMPLES:
  # Task started
  clog -type=task -state=in_progress -message="Adding VAT breakdown" -task-num="3/15" -session="nye-api"

  # Task completed
  clog -type=task -state=completed -message="VAT breakdown added" -task-num="3/15" -session="nye-api"

  # Question/blocked
  clog -type=question -state=blocked -message="Should VAT be inclusive or exclusive?" -session="nye-api"

  # Progress update
  clog -type=progress -message="50% complete (5/10 tasks)" -session="nye-api"

  # Session events
  clog -type=session -message="Started: API improvements design doc" -session="nye-api"

SUBJECTS (hardwired):
  task     -> claude.tasks
  question -> claude.questions
  progress -> claude.progress
  session  -> claude.session

EXIT CODES:
  0 - Success
  1 - Invalid arguments
  2 - NATS connection failed`)
}
