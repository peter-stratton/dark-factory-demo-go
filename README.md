# Dark Factory Demo: Go Bookmarks API

A simple Go REST API for managing bookmarks, designed as a hands-on tutorial for [Dark Factory](https://godarkfactory.com).

## What you'll learn

This demo walks you through real-world Dark Factory workflows using a Go project:

1. **Basic PR workflow** — Dark Factory creates a branch, opens a PR, and gets it merged
2. **Manual vs auto merge** — when to require human approval vs let changes land automatically
3. **Watch mode** — monitoring PRs and responding to human feedback
4. **Mechanical verification: tests** — using `go test` as a merge gate
5. **Mechanical verification: build** — ensuring `go build` passes before merge
6. **CI pipeline integration** — GitHub Actions required status checks with Dark Factory

## Prerequisites

- [Go 1.22+](https://go.dev/dl/)
- [golangci-lint](https://golangci-lint.run/welcome/install/)
- [Dark Factory CLI](https://godarkfactory.com)
- [GitHub CLI (`gh`)](https://cli.github.com/) — installed and authenticated
- A GitHub account
- An active [Anthropic](https://console.anthropic.com/) subscription with either an API key (`ANTHROPIC_API_KEY`) or OAuth token configured

---

## Getting started

### Step 1: Create your repo from this template

1. Click the **"Use this template"** button at the top of this page
2. Select **"Create a new repository"**
3. Name it whatever you like (e.g. `dark-factory-demo-go`) — public or private, either works
4. Click **"Create repository"**

### Step 2: Clone and verify

```bash
git clone https://github.com/YOUR_USERNAME/dark-factory-demo-go.git
cd dark-factory-demo-go
```

Make sure the project builds and tests pass:

```bash
make test    # run tests
make lint    # run linter
make build   # compile the binary
```

Optionally, start the server and poke around:

```bash
make run     # starts on :8080

# In another terminal:
curl http://localhost:8080/health
curl http://localhost:8080/bookmarks
```

### Step 3: Initialize Dark Factory

From the project root, run:

```bash
godark init --repo YOUR_USERNAME/dark-factory-demo-go
```

This sets up everything Dark Factory needs to work with your project. You should re-run `godark init` every time you update godark — it overwrites the default harness files and prompts with the latest versions but leaves your `godark.yaml` configuration and any custom prompts untouched.

- **`godark.yaml`** — project configuration (repo, build/test/lint commands, merge behavior)
- **`CLAUDE.md`** — the harness file that tells Claude Code about your project
- **`.claude/skills/`** — Claude Code skill definitions for planning and architecture
- **`docs/`** — architecture, conventions, and roadmap templates
- **`prompts/`** — agent prompt templates

Take a look at the generated `godark.yaml`:

```yaml
repo: "YOUR_USERNAME/dark-factory-demo-go"

auto_merge:
  feature: none       # human reviews every PR (we'll change this later)
  rollup: none

runtime:
  name: go
  version: "1.25.0"

# Build, test, format, and lint commands (auto-detected — verify for your project)
build_command: "go build ./..."
test_command: "go test ./..."
format_command: "gofmt -l -d ."
lint_command: "go vet ./..."
```

Dark Factory detected `go.mod` and filled in the runtime, build, test, format, and lint commands automatically. The `runtime` block tells the Docker sandbox which toolchain to install — without it, commands like `gofmt` and `go test` won't be available inside the container. Review the detected values and adjust if your project uses a different build chain (e.g. `make build` instead of `go build ./...`).

### Step 4: Run the doctor

Verify your environment is ready:

```bash
godark doctor
```

This checks the host-level requirements:

- Docker daemon running
- `gh` CLI installed and authenticated
- `ANTHROPIC_API_KEY` or `CLAUDE_CODE_OAUTH_TOKEN` set

Fix anything that fails before continuing. Language toolchains (Go, Python, etc.) are installed inside the Docker sandbox automatically — you don't need them on your host machine.

### Step 5: Commit the harness files

```bash
git add -A
git commit -m "Initialize Dark Factory harness"
git push
```

---

## Tutorial: Your first Dark Factory run

### Create a GitHub issue

Dark Factory works from GitHub issues. Create one for a small feature:

```bash
gh issue create \
  --title "Add a DELETE /bookmarks/:id endpoint" \
  --body "Add a DELETE endpoint that removes a bookmark by ID.

Acceptance criteria:
- DELETE /bookmarks/{id} returns 204 on success
- Returns 404 if the bookmark doesn't exist
- Include a test for both cases"
```

Note the issue number that's returned (e.g. `#1`).

### Create a milestone and assign the issue

```bash
gh api repos/YOUR_USERNAME/dark-factory-demo-go/milestones \
  --method POST -f title="Demo Phase 1"

gh issue edit 1 --milestone "Demo Phase 1"
```

### Run Dark Factory

```bash
godark run --milestone "Demo Phase 1" --repo YOUR_USERNAME/dark-factory-demo-go
```

Dark Factory will:

1. Fetch all open issues in the milestone
2. Resolve any dependency ordering
3. For each issue: create a branch, implement the change, run tests/lint/build, open a PR
4. Run quality and functional review on the PR
5. Label the PR `godark:awaiting-human-review` (since `auto_merge.feature` is set to `none`)

Watch the TUI as it works. When it finishes, you'll have an open PR on your repo ready for review.

### Review the PR

Head to your repo on GitHub and look at the open PR. Review the code, check that CI passes, and when you're satisfied:

- **Approve** the PR to merge it
- **Request changes** if you want Dark Factory to iterate

If you approve, merge it manually (since auto-merge is off). If you request changes, Dark Factory can pick them up in watch mode (covered next).

---

## Watch mode

Watch mode is a long-lived process that monitors your PRs and responds to human feedback. Start it with:

```bash
godark watch --repo YOUR_USERNAME/dark-factory-demo-go
```

With watch running:

1. Open a PR that's labeled `godark:awaiting-human-review`
2. Submit a review with **"Changes requested"** and describe what you want fixed
3. Watch detects the review, invokes the agent to address your feedback, and re-labels the PR for your next review
4. When you're happy, **Approve** the PR — watch merges it automatically

Press `Ctrl+C` to stop watch mode.

---

## Auto merge

Once you're comfortable with the workflow, you can let Dark Factory merge PRs without waiting for human approval.

Edit `godark.yaml`:

```yaml
auto_merge:
  feature: all    # auto-merge after both reviewers approve
```

Now when you run:

```bash
godark run --milestone "Demo Phase 1" --repo YOUR_USERNAME/dark-factory-demo-go
```

Dark Factory will implement, review, and merge PRs end-to-end without stopping for human review.

For a middle ground, use `low_risk` — Dark Factory auto-merges small, safe changes and flags larger ones for human review:

```yaml
auto_merge:
  feature: low_risk
```

---

## Mechanical verification

Dark Factory runs your project's build, test, and lint commands as merge gates. If any check fails, the agent fixes the issue and retries (up to `--max-retries`, default 3).

The commands come from your `godark.yaml`:

```yaml
build_command: "go build ./..."    # must compile
test_command: "go test ./..."      # tests must pass
lint_command: "golangci-lint run ./..." # no lint violations
```

The CI pipeline in `.github/workflows/ci.yml` runs the same checks on GitHub Actions. Dark Factory waits for required status checks to pass before merging.

---

## CI pipeline integration

This repo includes a GitHub Actions workflow that runs tests, linting, and build on every PR. To make Dark Factory respect these checks:

1. Go to your repo's **Settings > Branches**
2. Add a branch protection rule for `main`
3. Enable **"Require status checks to pass before merging"**
4. Select the `test`, `lint`, and `build` checks

Now Dark Factory won't merge a PR until CI is green, even with `auto_merge.feature: all`.

---

## Adding more issues

Try creating a few more issues to see Dark Factory handle a full milestone:

```bash
gh issue create \
  --title "Add a PUT /bookmarks/:id endpoint" \
  --body "Add a PUT endpoint that updates an existing bookmark.

Acceptance criteria:
- PUT /bookmarks/{id} accepts JSON body with url and/or title
- Returns the updated bookmark as JSON
- Returns 404 if the bookmark doesn't exist
- Returns 400 for invalid input
- Include tests for success, not-found, and validation cases"

gh issue create \
  --title "Add pagination to GET /bookmarks" \
  --body "Add query parameter pagination to the list endpoint.

Acceptance criteria:
- Support ?page=N&per_page=N query parameters
- Default to page=1, per_page=20
- Return a JSON envelope with items, page, per_page, and total fields
- Include tests for pagination behavior

Depends on: #1"
```

Assign them to your milestone and run again. Dark Factory resolves the `Depends on: #1` declaration and processes issues in the right order.

---

## Dry run

Not ready to let the agents loose? Use `--dry-run` to see what Dark Factory would do without making any changes:

```bash
godark run --milestone "Demo Phase 1" --dry-run --repo YOUR_USERNAME/dark-factory-demo-go
```

---

## What's next

- **Explore the harness** — edit `CLAUDE.md`, `docs/architecture.md`, and `docs/conventions.md` to shape how agents understand your project
- **Try `godark vet`** — validate your issues, architecture, and scenarios before running
- **Check analytics** — run `godark analyze` or `godark status` to see run history and metrics
- **Read the docs** — visit [godarkfactory.com](https://godarkfactory.com) for the full reference

---

## Project structure

```
.
├── main.go                        # Entry point
├── internal/
│   ├── model/
│   │   ├── bookmark.go            # Data types and validation
│   │   └── errors.go              # Domain errors
│   ├── server/
│   │   ├── server.go              # HTTP handlers and routing
│   │   └── server_test.go         # Handler tests
│   └── store/
│       ├── store.go               # In-memory data store
│       └── store_test.go          # Store tests
├── .github/workflows/ci.yml      # GitHub Actions CI
├── .golangci.yml                  # Linter configuration
├── Makefile                       # Build targets
└── go.mod                         # Go module definition
```

## Configuring the Claude model

Dark Factory uses whatever model is configured in your Claude Code project settings. To change the model (e.g. to use Sonnet for faster/cheaper runs, or Opus for more complex tasks), update your Claude Code settings — Dark Factory will pick it up automatically.
