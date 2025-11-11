# clog Contextual Reminders

This document shows what reminders appear for different event types and states.

## Task Events

### Task Started (with user-prompt)
```bash
$ clog -type=task -state=in_progress -user-prompt="Fix the login bug" -message="Debugging auth flow" -session="bugfix"
200 OK

  Remember: Log completion when done: clog -type=task -state=completed -message="..." -session="..."
```

### Task Started (without user-prompt)
```bash
$ clog -type=task -state=in_progress -message="Debugging auth flow" -session="bugfix"
200 OK

  TIP: Add -user-prompt="<exact user text>" to capture full context (VERBATIM user input)
```

### Task Completed
```bash
$ clog -type=task -state=completed -message="Login bug fixed" -session="bugfix"
200 OK

  Next: Log next task or use -type=progress for multi-step work
```

### Task Blocked
```bash
$ clog -type=task -state=blocked -message="Waiting for API key" -session="bugfix"
200 OK

  Consider: Use -type=question -state=blocked if waiting for user input
```

## Question Events

### Question (blocked)
```bash
$ clog -type=question -state=blocked -message="Should I use OAuth2 or SAML?" -session="bugfix"
200 OK

  Good: Always log questions before asking user
```

## Progress Events

### Progress Update
```bash
$ clog -type=progress -message="50% complete (5/10 tasks)" -session="bugfix"
200 OK
```
(No contextual reminder - progress updates are informational)

## Session Events

### Session Started
```bash
$ clog -type=session -message="Starting API improvements" -session="api-work"
200 OK

  Remember: Log tasks with -type=task as you work
```

### Session Completed
```bash
$ clog -type=session -state=completed -message="API improvements complete" -session="api-work"
200 OK
```
(No contextual reminder - session is ending)

## Design Rationale

**Character economy**: Every character costs money in LLM context windows.

**Contextual**: Reminders only appear when relevant to the current event type/state.

**Actionable**: Each reminder shows concrete next steps.

**Positive reinforcement**: "Good:" prefix confirms correct behavior (e.g., logging questions).

**No emojis**: Removed to save characters without losing clarity.
