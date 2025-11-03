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
- Test changes before committing

## Project Context

This is a Wingspan board game scoring application built with Go, featuring:
- SQLite database for game results storage
- Web-based UI for score tracking
- Podman (Docker) containerization support
    - Deployed to Kubernetes in production.