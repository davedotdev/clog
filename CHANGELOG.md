# clog Changelog

## v0.2.0 - Enhanced Feedback Loop & Configurable Reminders

### Changes
- **Context-aware feedback system**: clog now provides intelligent, context-specific reminders based on event type and state
- **User-configurable reminders**: Three customizable reminder slots collected at build time and baked into the binary
- **Reduced feedback noise**: Reminders are only shown when relevant to the current action
- **Smart reminder logic**:
  - Task start reminds to log completion
  - Task completion suggests next steps
  - Missing user-prompt triggers a helpful tip
  - Question logging provides positive reinforcement
- **printReminders function**: New architecture separates user-configured and context-specific reminders
- **getContextReminder function**: Implements intelligent reminder selection based on workflow state

### Key Feedback Improvements:
- **For tasks (in_progress)**:
  - If no user-prompt: "TIP: Add -user-prompt=\"<exact user text>\" to capture full context"
  - If user-prompt provided: "Remember: Log completion when done"
- **For tasks (completed)**: "Next: Log next task or use -type=progress for multi-step work"
- **For tasks (blocked)**: "Consider: Use -type=question -state=blocked if waiting for user input"
- **For questions (blocked)**: "Good: Always log questions before asking user"
- **For sessions**: "Remember: Log tasks with -type=task as you work"

### Custom Reminders:
Users can now add up to 3 custom reminders during `make build`:
- Examples: "Remember to ask me before deploying to production"
- Examples: "Always run tests before committing code"
- Examples: "Check the CI/CD pipeline status before proceeding"

These reminders appear in every clog output, keeping important project-specific context visible in the conversation history.

### Technical Details:
- Added `reminder1`, `reminder2`, `reminder3` variables in cmd/main.go (lines 36-38)
- Updated build process to collect and inject reminders
- Implemented `printReminders()` and `getContextReminder()` functions
- Version bumped to 0.2.0

### Why This Update?
Creates a comprehensive feedback loop that helps AI agents maintain context over long conversations. The combination of custom reminders and intelligent context-specific feedback reduces errors and keeps workflows on track.

---

## v0.1.4 - Mandatory Question Logging

### Changes
- **Updated CLAUDE.md**: Made question logging explicitly MANDATORY before asking any question
- **Added "When to log questions" section**: Clear list of when questions must be logged
- **Added "Common Mistakes" section**: Shows wrong vs correct examples
- **Emphasized in examples**: Added question logging example in workflow
- **No code changes**: Documentation-only update to v0.1.3

### Why This Update?
Claude was forgetting to log questions before asking them, especially when asking for clarification, preferences, or approach choices. This update makes it crystal clear that ALL questions must be logged first, with no exceptions.

### Key Rules Added:
- **ALWAYS log BEFORE asking ANY question**
- If you're about to ask, log it first
- No exceptions!

### Question Logging Triggers:
- Before asking for clarification
- Before asking which approach to take
- Before asking about preferences or choices
- Before asking for approval or confirmation
- Before asking about missing information
- ANY time you need user input to proceed

---

## v0.1.3 - Context Reinforcement Loop

### Changes
- **Added context reinforcement in output**: clog now provides feedback in its output to reinforce proper usage
- **Positive reinforcement**: When `-user-prompt` is present on task events, shows "✓ User prompt captured - context preserved for logging"
- **Missing prompt reminder**: When `-user-prompt` is missing on task events, shows warning with tip to add verbatim user input
- **Updated `printSuccess` function**: Now accepts and checks `userPrompt` parameter to provide contextual feedback

### Why This Update?
Creates a feedback loop in Claude's context window. Every clog output now includes a reminder about using `-user-prompt`, which helps counter context drift and prevents forgetting to use clog properly over long conversations.

### Example Output:

**With user-prompt (positive reinforcement):**
```
✓ Message published successfully to 'claude.tasks.started'
  Type: task
  State: in_progress
  Session: session-id

  ✓ User prompt captured - context preserved for logging
```

**Without user-prompt (gentle reminder):**
```
✓ Message published successfully to 'claude.tasks.started'
  Type: task
  State: in_progress
  Session: session-id

  ⚠️  TIP: Add -user-prompt="<exact user text>" to capture full context
       Use VERBATIM user input - don't summarize or paraphrase
```

---

## v0.1.2 - Verbatim Prompt Emphasis

### Changes
- **Documentation clarification**: Emphasized that `-user-prompt` must contain EXACT, VERBATIM user input
- **Updated CLAUDE.md**: Added explicit warnings against summarizing, paraphrasing, shortening, or rephrasing
- **Updated help text**: Enhanced with clear examples showing correct vs incorrect usage
- **Key change**: This is a documentation-only release - no code logic changes from v0.1.1

### Why This Update?
The v0.1.1 release added the `-user-prompt` flag but didn't sufficiently emphasize that it should capture the user's EXACT text. This led to Claude summarizing prompts instead of capturing them verbatim, defeating the purpose of the field.

### Critical Rule Added:
**ALWAYS** copy the user's prompt EXACTLY as they wrote it. Do NOT:
- Summarize it
- Paraphrase it
- Shorten it
- Rephrase it

---

## v0.1.1 - User Prompt Logging

# clog v0.1.1 - Change Summary

## Overview
Added user prompt logging capability to allow Claude to log both what the user requested and what action is being taken, in a single clog call.

## Changes Made

### 1. Go Source Code (`cmd/main.go`)

#### Added UserPrompt Field to Message Struct (line 51)
```go
type Message struct {
	Event      string `json:"event"`
	Timestamp  string `json:"timestamp"`
	SessionID  string `json:"session_id,omitempty"`
	Message    string `json:"message"`
	UserPrompt string `json:"user_prompt,omitempty"`  // NEW FIELD
	State      string `json:"state,omitempty"`
	TaskNum    string `json:"task_num,omitempty"`
}
```

#### Added -user-prompt Flag (line 65)
```go
userPromptFlag := flag.String("user-prompt", "", "User's input prompt (optional)")
```

#### Updated Message Creation (line 102)
```go
msg := Message{
	Event:      subject,
	Timestamp:  time.Now().UTC().Format(time.RFC3339),
	SessionID:  *sessionFlag,
	Message:    *messageFlag,
	UserPrompt: *userPromptFlag,  // NEW FIELD
	State:      *stateFlag,
	TaskNum:    *taskNumFlag,
}
```

#### Enhanced Help Text (lines 306-334)
- Added `-user-prompt` to OPTIONAL FLAGS section
- Created prominent "CRITICAL WORKFLOW" section with visual separators
- Documented the pattern: use both `-user-prompt` and `-message` in ONE call
- Updated examples to show new usage

### 2. Global Instructions (`~/.claude/CLAUDE.md`)

#### Added Prominent Banner and Workflow
- Visual separator with `═` characters at top of file
- "CRITICAL: ALWAYS USE CLOG AT THE START OF EVERY TASK-BASED RESPONSE"
- MANDATORY WORKFLOW section with clear STEP 1 and STEP 2
- Emphasized: `-user-prompt` = what user asked, `-message` = what you're doing
- Reorganized entire file to prioritize the critical workflow

## Key Concepts

### The Two Fields:
- **`-user-prompt`**: What the user asked for (their exact request)
- **`-message`**: What Claude is doing about it (the action being taken)

### Usage Pattern:
```bash
# User says: "Add logging to the API endpoints"
clog -type=task -state=in_progress \
  -user-prompt="Add logging to the API endpoints" \
  -message="Adding structured logging middleware" \
  -session="api-logging"
```

### Benefits:
- One clog call captures full context (request + action)
- No extra noise for the user
- Clean NATS payload with both user intent and Claude's response

## Build Instructions

To rebuild and install:

```bash
# From the clog directory
make build

# This will:
# 1. Backup main.go
# 2. Inject NATS configuration
# 3. Build the binary
# 4. Restore main.go to template state
# 5. Output: ./clog binary

# Then install globally (if needed):
sudo cp ./clog /usr/local/bin/clog
sudo chmod +x /usr/local/bin/clog
```

## Testing

After rebuild, verify with:
```bash
clog -h | grep -A 2 "user-prompt"
```

Should show:
```
  -user-prompt User's input prompt (what the user asked for)
```

## JSON Payload Example

Before (v0.1.0):
```json
{
  "event": "claude.tasks.started",
  "timestamp": "2025-10-15T12:00:00Z",
  "session_id": "api-work",
  "message": "Adding logging middleware",
  "state": "in_progress"
}
```

After (v0.1.1):
```json
{
  "event": "claude.tasks.started",
  "timestamp": "2025-10-15T12:00:00Z",
  "session_id": "api-work",
  "message": "Adding logging middleware",
  "user_prompt": "Add logging to the API endpoints",
  "state": "in_progress"
}
```

## Version

Version bumped to: **0.1.1** (line 16 in cmd/main.go)
