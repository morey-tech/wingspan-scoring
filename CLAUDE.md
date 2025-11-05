# Wingspan Scoring - Development Guidelines

## Git Workflow

- **Always create a new feature branch before making any code changes**
- Never commit directly to `main` branch
- Branch naming: use descriptive names (e.g., `feature/score-calculator`, `fix/round-goals-bug`)

## Commit Standards

- Use conventional commit format: `type: description`
- Types: `feat`, `fix`, `refactor`, `docs`, `test`, `chore`
- Keep commits focused and atomic
- Write clear, descriptive commit messages
- Don't include Claude Code in the commit message.

## Code Standards

- Follow Go best practices and idioms
- Run `go fmt` before committing
- Ensure code builds successfully with `go build`
- Run tests and ensure they pass before committing
- Write tests for new functionality
- Maintain test coverage above 70%

## Testing

### Running Tests

Run all tests:
```bash
go test ./...
```

Run tests with verbose output:
```bash
go test -v ./...
```

Run tests with coverage:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

View coverage in browser:
```bash
go tool cover -html=coverage.out
```

Run tests with race detection:
```bash
go test -race ./...
```

### Test Structure

Tests are organized alongside source files:
- `goals/scorer_test.go` - Round goal scoring tests
- `goals/selector_test.go` - Goal selection tests
- `goals/goals_test.go` - Goal definition tests
- `scoring/scoring_test.go` - End-game calculation tests
- `db/db_test.go` - Database initialization tests
- `db/game_results_test.go` - CRUD operation tests
- `main_test.go` - HTTP handler tests

### Writing Tests

- Use the `testify/assert` library for assertions
- Name test functions with the pattern `Test<FunctionName>_<Scenario>`
- Use table-driven tests for multiple scenarios
- Test both success and error cases
- Use temporary databases for database tests

### CI/CD Integration

Tests run automatically on every push to main via GitHub Actions. The build will fail if:
- Any test fails
- Race conditions are detected
- Code doesn't compile

## Project Context

This is a Wingspan board game scoring application built with Go, featuring:
- SQLite database for game results storage
- Web-based UI for score tracking
- Podman (Docker) containerization support
    - Deployed to Kubernetes in production.