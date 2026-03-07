# Contributing to Contributum

First off, thanks for your interest in contributing! Contributum is an open-source project and we welcome contributions of all kinds.

## Getting Started

### Prerequisites

- **Go 1.24 or later** — Check your version with `go version`
- **Git** — For cloning and making commits
- **GitHub Personal Access Token** (optional but recommended)
  - Provides higher rate limits for the GitHub API (30 req/min vs 10 req/min)
  - Create one at https://github.com/settings/tokens (no special scopes needed for public repo search)

### Clone and Build

```bash
# Clone the repository
git clone https://github.com/ijshd7/contributum.git
cd contributum

# Install dependencies
go mod download

# Build the binary
go build -o contributum .

# Test the build
./contributum search --lang go --topic cli --skill beginner
```

### Authentication (Optional)

```bash
# Copy the example environment file
cp .env.example .env

# Edit .env and add your GitHub token
# Or set it directly in your shell:
export GITHUB_TOKEN=your_token_here
```

## Development Workflow

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Run tests for a specific package
go test ./internal/scoring
```

### Code Style

We follow standard Go conventions:

- Run `gofmt` on your code: `go fmt ./...`
- Use meaningful variable names and comments for non-obvious logic
- Keep functions focused and testable
- Follow [Effective Go](https://golang.org/doc/effective_go) principles

### Commit Message Format

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Type:**
- `feat` — A new feature
- `fix` — A bug fix
- `test` — Adding or improving tests
- `docs` — Documentation changes
- `chore` — Build, dependencies, etc.
- `refactor` — Code refactoring without feature changes

**Scope** (optional): The module or package affected (e.g., `ghapi`, `scoring`, `ui`)

**Subject:** Concise imperative description (e.g., "add repo search by language")

**Body** (optional): Detailed explanation of the change and why it's needed.

**Examples:**
```
feat(ghapi): add support for filtering by topic

fix(scoring): exclude repos with 0 good first issues for beginners

test(ui): add tests for truncate helper function
```

## Pull Request Process

1. **Fork the repository** and create a feature branch from `master`:
   ```bash
   git checkout -b feat/your-feature-name
   ```

2. **Make your changes** and commit with descriptive messages.

3. **Write or update tests** for your changes. Ensure all tests pass:
   ```bash
   go test ./...
   ```

4. **Push to your fork**:
   ```bash
   git push origin feat/your-feature-name
   ```

5. **Open a pull request** to the `master` branch with a clear description of your changes.

6. **Address feedback** from reviewers — we'll work with you to refine the changes.

## What We're Looking For

- **Bug fixes** — Especially with tests demonstrating the fix
- **New features** — With tests and documentation
- **Improved tests** — Higher coverage is always welcome
- **Documentation** — Typos, clarity, or new guides
- **Performance improvements** — With benchmarks showing the impact
- **Code quality** — Refactoring or simplification of complex logic

## Project Structure

```
contributum/
├── cmd/                    # CLI commands (search, issues, version)
├── internal/
│   ├── ghapi/             # GitHub API client and search logic
│   ├── scoring/           # Repository scoring algorithm
│   └── ui/                # Terminal UI and output rendering
├── main.go                # Program entry point
├── README.md              # User documentation
└── CONTRIBUTING.md        # This file
```

### Key Packages

- **`cmd/`** — Cobra command definitions. Entry points for `search`, `issues`, and `version` commands.
- **`internal/ghapi/`** — GitHub API client wrapper and search query builders. Handles authentication and result enrichment.
- **`internal/scoring/`** — Scoring algorithm that ranks repositories by activity, friendliness, and relevance.
- **`internal/ui/`** — Terminal rendering using Lipgloss and Bubbletea. Handles table formatting and JSON output.

## Testing Patterns

### Unit Testing Pure Functions

For pure functions (no I/O, no side effects), write straightforward unit tests:

```go
func TestTruncate(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"hello", 10, "hello"},              // no truncation
		{"hello world", 8, "hello …"},       // truncated with ellipsis
		{"hi", 2, "h…"},                     // exact boundary
	}

	for _, tt := range tests {
		if got := truncate(tt.input, tt.maxLen); got != tt.expected {
			t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.expected)
		}
	}
}
```

### Testing with io.Writer

For functions that write output, pass a `bytes.Buffer` to capture and verify output:

```go
func TestRenderRepoTable(t *testing.T) {
	var buf bytes.Buffer
	repos := []scoring.ScoredRepo{/* ... */}

	renderRepoTable(repos, &buf)

	output := buf.String()
	if !strings.Contains(output, "Repository Search Results") {
		t.Error("expected title in output")
	}
}
```

## Need Help?

- **Questions?** Open a [GitHub Discussion](https://github.com/ijshd7/contributum/discussions)
- **Found a bug?** Open a [GitHub Issue](https://github.com/ijshd7/contributum/issues)
- **Want to chat?** Discussions are a great place to brainstorm features or get feedback

## Code of Conduct

Please note that this project is released with a [Contributor Code of Conduct](CODE_OF_CONDUCT.md). By participating in this project you agree to abide by its terms.

---

Happy contributing! 🎉
