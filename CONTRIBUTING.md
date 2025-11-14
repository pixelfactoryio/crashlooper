# Contributing to Crashlooper

Thank you for your interest in contributing to Crashlooper! This guide will help you understand our development workflow and commit conventions.

## Table of Contents

- [Commit Message Convention](#commit-message-convention)
- [Semantic Versioning](#semantic-versioning)
- [Development Workflow](#development-workflow)
- [Making Changes](#making-changes)
- [Creating a Release](#creating-a-release)

## Commit Message Convention

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification for our commit messages. This allows us to automatically determine version bumps and generate changelogs.

### Format

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

### Types

The commit type determines how the version number will be bumped:

| Type | Description | Version Bump | Example |
|------|-------------|--------------|---------|
| `feat` | A new feature | MINOR (0.1.0 → 0.2.0) | `feat: add memory increment feature` |
| `fix` | A bug fix | PATCH (0.1.0 → 0.1.1) | `fix: correct memory leak in handler` |
| `docs` | Documentation only changes | None | `docs: update installation instructions` |
| `style` | Code style changes (formatting, etc.) | None | `style: fix gofmt formatting` |
| `refactor` | Code refactoring without feature changes | None | `refactor: simplify crash handler logic` |
| `perf` | Performance improvements | PATCH (0.1.0 → 0.1.1) | `perf: optimize memory allocation` |
| `test` | Adding or updating tests | None | `test: add unit tests for logger` |
| `chore` | Build process or auxiliary tool changes | None | `chore: upgrade Go to 1.25` |
| `ci` | CI/CD configuration changes | None | `ci: add automated testing workflow` |
| `build` | Changes that affect the build system | None | `build: update dependencies` |

### Breaking Changes

To indicate a BREAKING CHANGE (which triggers a MAJOR version bump), use one of these methods:

1. **Add `!` after the type:**
   ```
   feat!: remove deprecated --legacy flag
   ```

2. **Add `BREAKING CHANGE:` in the footer:**
   ```
   feat: redesign configuration API

   BREAKING CHANGE: the configuration file format has changed from JSON to YAML
   ```

Breaking changes trigger: `1.2.3 → 2.0.0`

### Scope (Optional)

The scope provides additional context about what part of the codebase is affected:

```
feat(api): add new endpoint for health checks
fix(docker): correct entrypoint path
docs(readme): add usage examples
```

Common scopes for this project:
- `api` - HTTP API changes
- `docker` - Docker-related changes
- `cli` - Command-line interface
- `build` - Build system and tooling
- `ci` - Continuous integration

### Description

- Use imperative, present tense: "add" not "added" or "adds"
- Don't capitalize the first letter
- No period (.) at the end
- Keep it concise (50 characters or less is ideal)

### Examples

**Good commits:**
```
feat: add configurable timeout support
fix: correct memory leak in crash handler
docs: update docker deployment guide
test: add integration tests for HTTP endpoints
chore: upgrade dependencies to latest versions
ci: add automated tagging workflow
```

**Bad commits:**
```
Updated stuff                    # Too vague, no type
Fix bug.                         # Not descriptive enough
Added new feature for users      # Past tense, should be present
FEAT: New API endpoint           # Type should be lowercase
fix: Fixed the crash bug.        # Mixed tense, unnecessary period
```

### Multi-line Commits

For more complex changes, add a body and/or footer:

```
feat: add crash recovery mechanism

Implements automatic restart capability when the process crashes
unexpectedly. The recovery mechanism includes exponential backoff
and configurable retry limits.

Closes #42
```

## Semantic Versioning

This project follows [Semantic Versioning](https://semver.org/) (SemVer):

```
MAJOR.MINOR.PATCH (e.g., v1.2.3)
```

- **MAJOR** (1.x.x): Breaking changes that are not backward compatible
- **MINOR** (x.1.x): New features that are backward compatible
- **PATCH** (x.x.1): Bug fixes and minor improvements

### Version Bump Rules

Based on your commits since the last tag:

| Commits | Version Bump | Example |
|---------|--------------|---------|
| Only `fix:`, `perf:` | PATCH | v1.0.0 → v1.0.1 |
| Contains `feat:` | MINOR | v1.0.0 → v1.1.0 |
| Contains `BREAKING CHANGE` or `!` | MAJOR | v1.0.0 → v2.0.0 |
| Only `docs:`, `chore:`, `ci:`, etc. | None | No new tag |

## Development Workflow

### 1. Check Current Version

Before making changes, check the current version:

```bash
make current-version
```

### 2. Make Your Changes

Create a branch and make your changes with proper conventional commits:

```bash
git checkout -b feature/my-feature
# Make changes
git add .
git commit -m "feat: add my awesome feature"
```

### 3. Check Next Version

Before creating a PR, verify what version will be generated:

```bash
make next-version
```

This shows:
- What the next version will be
- All commits since the last tag
- Whether a version bump is warranted

### 4. Create Pull Request

Push your branch and create a pull request. The CI will automatically:
- Run tests
- Build Docker images with tag `pr-<number>-<sha>`
- Validate your changes

### 5. Merge to Main

Once approved and merged to `main`, the auto-tagging workflow will:
- Analyze all commits since the last tag
- Determine the appropriate version bump
- Create and push the new tag automatically
- Trigger the release workflow to build and publish artifacts

## Creating a Release

### Automatic (Recommended)

Simply merge your PR to `main`. The GitHub Actions workflow will automatically:
1. Analyze commits using conventional commit format
2. Determine the next version based on commit types
3. Create and push the appropriate tag
4. Trigger the release build

### Manual

If you need to create a release manually:

#### Smart Release (Based on Commits)

```bash
make release
```

This will:
- Analyze commits since the last tag
- Determine the next version
- Ask for confirmation
- Create and push the tag

#### Force Specific Version Bump

If you need to override the automatic detection:

```bash
make release-patch   # Force patch bump (0.1.0 → 0.1.1)
make release-minor   # Force minor bump (0.1.0 → 0.2.0)
make release-major   # Force major bump (0.1.0 → 1.0.0)
```

## Testing Your Changes

Before committing, ensure all tests pass:

```bash
# Run unit tests
make test

# Run linter
make lint

# Run vet
make vet
```

## Questions?

If you have questions about contributing or the commit convention, please:
1. Check this guide first
2. Review existing commits for examples
3. Open an issue for clarification

Thank you for contributing to Crashlooper!
