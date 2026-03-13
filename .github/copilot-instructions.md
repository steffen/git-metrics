# Copilot Instructions for git-metrics

## Solution Approach

- Prefer the simplest solution that meets requirements.
- Avoid over-engineering; implement only what's needed to solve the immediate problem.
- Keep dependencies and abstractions to a minimum.

## Naming Conventions

- Do not use abbreviations in source code (variable names, function names, constants, etc.).
- Use full, descriptive words (e.g., `initialBufferSize` not `initialBufSize`, `configuration` not `config`).
- Exception: well-established Go idioms (`err`, `ctx`, `fmt`, `i` in loops) are acceptable.

## Git Workflow

- Make git commits only through the git CLI. Do not use the GitHub API for creating commits.

## Git Commits

- **Title**: Use a concise imperative sentence describing *what* the commit does (e.g., `Fix reader buffer to prevent token too long errors`, not `Fixed stuff` or `Update reader.go`).
- **Description**: Explain *why* the change is needed and *what* it does at a high level — not a line-by-line diff summary.
- Separate title from description with a blank line.
- In this project, keep the commit title under 72 characters (our chosen limit) and wrap body lines at 72 characters.
- Reference related issues or PRs when applicable (e.g., `Fixes #76`).
- Prefer atomic commits — one logical change per commit; use `git add -p` or stage files selectively when multiple unrelated changes are in the working tree.

## Pull Requests

- Keep the PR description accurate and in sync with the actual diff — after undoing commits or implementing Copilot suggestions, update the description to reflect the real changes.
- When resolving Copilot review suggestions: implement the fix, reply to the thread with the commit SHA and a one-line summary of what was done, then resolve the thread.
- When dismissing a Copilot suggestion (e.g., one that contradicts project guidelines), reply explaining the design decision before resolving.

## Agent Workflow

- When creating intermediate or working files during agent or chat processes, store them in the `tmp/` folder at the workspace root, not in the host machine's `/tmp` directory. This avoids file access errors and keeps working artifacts within the project.

## Code Style

- Use tabs for indentation in Go files.
- In Go files, empty lines should not contain indentation characters.

## Go Idioms

- Use sentinel error values instead of string comparison (e.g., `errors.Is(err, bufio.ErrTooLong)` not `err.Error() == "bufio.Scanner: token too long"`).
- Hoist `regexp.MustCompile` calls to package-level `var` declarations so regexes are compiled once and reused, not on every function call.
- All warning/error output goes to `os.Stderr`; report output goes to `os.Stdout` — never mix them.
- When multiple packages need the same logic, extract it to `pkg/utils` as an exported function rather than duplicating it — duplicated logic diverges silently.

## Error Handling & CLI

- Never use interactive prompts that require user input during execution.
- When errors occur, provide clear, non-blocking error messages with actionable guidance.
- Validate user-provided flag values explicitly and exit with an actionable error rather than silently adjusting them.

## Project Structure

- The project is a Go CLI tool that analyzes git repositories and generates metrics reports.
- The main entry point is `main.go` at the project root.
- Packages are organized under `pkg/` with the following structure:
  - `pkg/display/` — output formatting and display logic
  - `pkg/display/sections/` — individual report section renderers
  - `pkg/git/` — git command execution and data collection
  - `pkg/models/` — shared data models and types
  - `pkg/progress/` — progress indicators and spinner animations
  - `pkg/requirements/` — system requirement checks
  - `pkg/utils/` — utility functions (formatting, terminal detection, system info)
- Test fixtures are stored in `fixtures/`.
- Build and test scripts are in `script/`.
- The `tmp/` directory is used for test repositories and working files.

## Testing

- Run unit tests with `go test ./pkg/...`.
- Run integration tests with `script/run-integration-tests`.
- Update test fixtures with `script/update-fixtures`.
- The `fixtures/git-metrics.txt` file contains platform-dependent values (e.g., object sizes differ between macOS and Linux).
- When changing stdout output (adding/removing lines, changing formatting), update all affected fixture files.
- Validate fixture changes locally by running `script/run-integration-tests` or update fixtures with `script/update-fixtures`.

## Output Conventions

- Report output goes to `os.Stdout`.
- Error messages go to `os.Stderr`.
- Progress indicators use ANSI escape sequences for terminal manipulation.
