You're taking over ticket `<TICKET-ID>` in repo `<REPO-PATH>`.

- Read docmgr basics first
  - Run:
    ```bash
    cd <REPO-PATH>

    docmgr help how-to-use
    docmgr list tickets --ticket <TICKET-ID>
    docmgr list docs --ticket <TICKET-ID>
    docmgr tasks list --ticket <TICKET-ID>
    ```

- Read up on this ticket (in order)
  1) The ticket's main index document (usually `index.md` in the ticket directory)
  2) Implementation diaries (if any):
     - Check the `log/` directory for diary files
     - Read them chronologically to understand the implementation history
  3) Tasks and changelog:
     - `tasks.md` - list of tasks to complete
     - `changelog.md` - record of changes made
  4) Background documentation:
     - Any relevant guides or tutorials referenced in the ticket
     - Check `docmgr list docs --ticket <TICKET-ID>` for related documentation

- Start tasks
  - List and begin with the first unchecked task:
    ```bash
    docmgr tasks list --ticket <TICKET-ID>
    ```
  - As you complete a task:
    ```bash
    docmgr tasks check --ticket <TICKET-ID> --id <ID>
    ```

- At EVERY step: relate files and keep a changelog
  - After creating/updating files:
    ```bash
    docmgr relate --ticket <TICKET-ID> \
      --file-note "/ABSOLUTE/PATH/HERE:Short note on why this file matters"

    docmgr changelog update --ticket <TICKET-ID> \
      --entry "What changed and why" \
      --file-note "/ABSOLUTE/PATH/HERE:Reason"
    ```

- Keep a running implementation diary
  - Append brief entries after meaningful steps (continue in the existing diaries or create a new `various/` note). Use:
    ```text
    Write up an implementation diary for what you just did, documenting every step, what you had to do, what worked, what didn't work, what you learned, what you should do in the future,
    ```

- Repo-specific information
  - Document any repo-specific build commands, test procedures, or comparison artifacts here
  - Note any special setup or configuration needed for this ticket
  - Include paths to important files or directories relevant to the work

- Known issues and immediate focus
  - Document any known gaps, blockers, or areas requiring attention
  - List outstanding tasks or follow-up items
  - Note any workarounds or temporary solutions in place

- Useful docmgr helpers
  ```bash
  docmgr status --summary-only
  docmgr list tickets
  docmgr meta update --ticket <TICKET-ID> --field Status --value active
  ```

- Where to start next (suggested)
  - Review the tasks list and identify the highest priority unchecked task
  - Check the ticket's index document for roadmap or next steps
  - Review any related tickets or dependencies mentioned in the documentation

