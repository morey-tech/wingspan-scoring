# Development Workflow

## Standard Issue Resolution Workflow

This workflow follows Test-Driven Development (TDD) principles and GitHub best practices for transparent issue tracking.

### Step 1: Write Test First (TDD)

Before creating the issue or implementing any fix, write a failing test that demonstrates the bug or validates the expected behavior.

**Why First?**
- Ensures complete understanding of the problem
- Provides concrete reproduction case
- Validates the fix works when test passes

**Example (Go backend test)**:
```bash
# Add test to appropriate test file
# Run to verify it fails
go test -v -run TestName
```

### Step 2: Create GitHub Issue

Document the bug/feature in the issue tracker with full context.

**Command**:
```bash
gh issue create \
  --title "Bug: descriptive title" \
  --body "$(cat <<'EOF'
## Description
Clear description of the bug or feature request

## Steps to Reproduce (for bugs)
1. Step one
2. Step two
3. Observed behavior

## Expected Behavior
What should happen

## Actual Behavior
What currently happens

## Technical Details
- Root cause (if known)
- Affected files
- Suggested fix approach

## Environment
- Version/branch
- Configuration details
EOF
)"
```

**Save the issue number** (e.g., `#92`) for later reference.

### Step 3: Create Feature Branch

Always work in a feature branch, never commit directly to `main`.

**Command**:
```bash
git checkout -b fix/descriptive-name
# or
git checkout -b feature/descriptive-name
```

**Naming Convention**:
- Bug fixes: `fix/issue-description`
- Features: `feature/issue-description`
- Refactors: `refactor/issue-description`

### Step 4: Implement Fix/Feature

Make the code changes to resolve the issue and make the test pass.

**Workflow**:
1. Implement the fix
2. Run tests to verify they pass
3. Verify all existing tests still pass
4. Test manually if applicable

**Commands**:
```bash
# Run specific test
go test -v -run TestName

# Run all tests
go test ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
```

### Step 5: Commit Without Claude Attribution

Create a clean, professional commit message following conventional commit format.

**Command**:
```bash
git add <files>

git commit -m "type: brief description

- Detailed change 1
- Detailed change 2
- Detailed change 3

Additional context or explanation of why the change was needed.

Closes #<issue-number>"
```

**Commit Types**:
- `feat`: New feature
- `fix`: Bug fix
- `refactor`: Code restructuring
- `docs`: Documentation changes
- `test`: Test additions/changes
- `chore`: Maintenance tasks

**Important**: Do NOT include "Generated with Claude Code" or similar attribution.

### Step 6: Push Branch

Push the feature branch to remote repository.

**Command**:
```bash
git push -u origin fix/descriptive-name
```

### Step 7: Create PR with Detailed Explanation

Create a pull request with comprehensive explanation and enable auto-squash merge.

**Command**:
```bash
gh pr create \
  --title "fix: brief description matching commit" \
  --body "$(cat <<'EOF'
## Summary
High-level overview of what this PR does

## Problem Statement
Detailed explanation of the bug or need for the feature

## Solution
How this PR addresses the problem:
- Key change 1
- Key change 2
- Key change 3

## Technical Details
### Root Cause (for bugs)
Explanation of what caused the issue

### Implementation Approach
How the fix works, including:
- Architecture decisions
- Code patterns used
- Any trade-offs made

### Code Changes
- `file1.go`: Description of changes
- `file2.js`: Description of changes

## Testing
- [x] Backend tests pass
- [x] Manual testing completed
- [x] Regression testing done

### Test Coverage
- New test: TestName - validates X
- Existing tests: All passing

### Manual Test Scenarios
1. Scenario 1: Expected result
2. Scenario 2: Expected result

## Side Effects
Analysis of any side effects or impacts:
- Performance impacts
- Breaking changes
- Migration requirements

## Related Issues
Closes #<issue-number>

## Additional Notes
Any other context reviewers should know
EOF
)" && \
gh pr merge --auto --squash
```

**The `--auto --squash` flags**:
- `--auto`: Enables auto-merge when CI checks pass
- `--squash`: Squashes all commits into one on merge
- Automatically closes linked issue on merge

### Step 8: Verify Auto-Merge Status

Check that auto-merge is enabled and CI is running.

**Commands**:
```bash
# Check PR status
gh pr status

# View PR checks
gh pr checks

# View PR details
gh pr view
```

---

## GitHub CLI Setup

### First-Time Authentication

Before using `gh` commands, authenticate with GitHub:

```bash
gh auth login
```

Follow the interactive prompts to authenticate via browser or token.

---

## Code Standards

### Go Best Practices

Before committing Go code, ensure:

```bash
# Format code
go fmt ./...

# Build successfully
go build

# All tests pass
go test ./...

# No race conditions
go test -race ./...
```

### Requirements
- Follow Go idioms and best practices
- Maintain test coverage above **70%**
- Run `go fmt` before every commit
- Ensure code builds without errors
- All tests must pass

---

## Testing Guide

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific test
go test -v -run TestName

# Run tests with coverage
go test -coverprofile=coverage.out ./...

# View coverage summary
go tool cover -func=coverage.out

# View coverage in browser
go tool cover -html=coverage.out

# Run with race detection
go test -race ./...
```

### Test Organization

Tests are organized alongside source files:
- `goals/scorer_test.go` - Round goal scoring tests
- `goals/selector_test.go` - Goal selection tests
- `goals/goals_test.go` - Goal definition tests
- `scoring/scoring_test.go` - End-game calculation tests
- `db/db_test.go` - Database initialization tests
- `db/game_results_test.go` - CRUD operation tests
- `main_test.go` - HTTP handler tests

### Test Writing Guidelines

- Use `testify/assert` library for assertions
- Name pattern: `Test<FunctionName>_<Scenario>`
- Use table-driven tests for multiple scenarios
- Test both success and error cases
- Use temporary databases for database tests

### CI/CD Integration

Tests run automatically on every push via GitHub Actions. The build fails if:
- Any test fails
- Race conditions are detected
- Code doesn't compile
- Test coverage drops below threshold

---

## Alternative PR Creation Methods

### Basic PR Creation (Opens Editor)
```bash
gh pr create
```

### PR with Inline Title and Body
```bash
gh pr create --title "feat: add score validation" --body "Implements input validation for bird scores"
```

### Create PR and Open in Browser
```bash
gh pr create --web
```

### Enable Auto-Merge on Existing PR
```bash
# Enable auto-merge with squash on specific PR number
gh pr merge 42 --auto --squash

# Enable auto-merge on current branch
gh pr merge --auto --squash
```

---

## Example: Bug Fix Workflow (Issue #92)

### Real-world example of this workflow in action:

**Issue**: Game End section showing incorrect round goal scores for tied players

**Step 1: Write Test**
```bash
# Added test to main_test.go
go test -v -run TestHandleCalculateScores_TwoPlayerTieRound2
# Test fails as expected ✓
```

**Step 2: Create Issue**
```bash
gh issue create \
  --title "Bug: Game End section shows incorrect round goal scores for tied players" \
  --body "[detailed bug report]"
# Created Issue #92
```

**Step 3: Create Branch**
```bash
git checkout -b fix/round-goal-tie-scoring
```

**Step 4: Implement Fix**
- Refactored `calculatePlayerRoundGoalScore()` in `static/js/app.js`
- Made function async to call API for tie resolution
- Updated 12+ calling functions to handle async/await

**Step 5: Commit**
```bash
git commit -m "fix: correct round goal scoring to handle ties

- Add backend integration tests for tie scenarios
- Refactor calculatePlayerRoundGoalScore() to use API
- Update function calls to handle async/await
- Add fallback logic for API failures

Fixes bug where tied players received raw scores instead
of averaged points in Game End section.

Closes #92"
```

**Step 6: Push**
```bash
git push -u origin fix/round-goal-tie-scoring
```

**Step 7: Create PR**
```bash
gh pr create \
  --title "fix: correct round goal scoring to handle ties" \
  --body "[comprehensive PR description with root cause analysis]" && \
gh pr merge --auto --squash
# Created PR #93 with auto-merge enabled
```

**Step 8: Verify**
```bash
gh pr checks
# All checks passing, PR auto-merged ✓
# Issue #92 automatically closed ✓
```

---

## Benefits of This Workflow

1. **Test First (TDD)**: Ensures bug is fully understood before coding
2. **Issue Tracking**: Creates transparent record of problems and solutions
3. **Feature Branches**: Keeps main branch stable
4. **Clean Commits**: Professional git history without AI attribution
5. **Detailed PRs**: Makes review process efficient and educational
6. **Auto-Merge**: Reduces manual overhead while maintaining quality gates
7. **Automatic Issue Closure**: Links code changes to issue tracking

---

## Quick Reference

### One-liner for PR with Auto-Merge
```bash
gh pr create --title "type: description" --body "detailed explanation

Closes #123" && gh pr merge --auto --squash
```

### Check Workflow Status
```bash
# View issue
gh issue view 92

# View PR status
gh pr status

# View PR details
gh pr view 93

# Check CI status
gh pr checks
```

---

## Project Context

This is a Wingspan board game scoring application built with Go, featuring:
- SQLite database for game results storage
- Web-based UI for score tracking
- Podman (Docker) containerization support
- Kubernetes deployment in production

### Technology Stack
- **Backend**: Go with standard library HTTP server
- **Database**: SQLite with modernc.org/sqlite driver
- **Frontend**: Vanilla JavaScript, HTML, CSS
- **Containerization**: Podman/Docker
- **Orchestration**: Kubernetes
- **CI/CD**: GitHub Actions
