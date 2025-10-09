# Design Spec

A binary baked with credentials for connecting to a NATS server (override with --url, --creds). 
Posts on subjects: 

  Task Events:
  - claude.tasks.started - When I start a new todo item
  - claude.tasks.completed - When I complete a todo item
  - claude.tasks.blocked - When I need your input

  Question Events:
  - claude.questions.asked - When I interrupt with a question
  - claude.questions.waiting - Reminder I'm blocked on your answer

  Progress Events:
  - claude.progress.update - Overall completion percentage
  - claude.session.started - When you give me a new design doc
  - claude.session.completed - When all todos are done

Uses stateless core NATS publish.

Create a Go binary with fields for baked in URL and credentials (seed, JWT), then provide ENV var inputs for NATS_URL, NATS_CREDS and enough switches to make the binary.

Output of help: "-h"
```text
  clog - Claude Log Publisher for NATS

  USAGE:
    clog -type=<event_type> -message="<text>" [options]

  REQUIRED FLAGS:
    -type        Event type: task|question|progress|session
    -message     Message content (string)

  OPTIONAL FLAGS:
    -state       Task state: pending|in_progress|blocked|completed
    -task-num    Current task number (e.g., "3/15")
    -session     Session identifier (any string)

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
    2 - NATS connection failed
```

# Claude -> Put your todo list here and mark as you go


