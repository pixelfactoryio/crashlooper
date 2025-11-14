# Crashlooper

A simple container that crash 💥 after a set amount of time ⏰

## Usage

```bash
$ docker build -t crashlooper .
$ docker run --rm -it crashlooper --help
Usage:
  crashlooper [flags]

Flags:
      --crash-after duration                 Server will crash itself after specified period (default=0 means never)
  -h, --help                                 help for crashlooper
      --log-level string                     Server log level (default "info")
      --memory-increment string              crashlooper memory usage increment
      --memory-increment-interval duration   crashlooper memory usage increment interval (default 1s)
      --memory-target string                 crashlooper memory usage target
      --port string                          Server bind port (default "3000")
```

## Example

```bash
docker run --rm -it pixelfactory/crashlooper:latest --crash-after 10s
```

## Docker Images

Pre-built Docker images are available on Docker Hub: `pixelfactory/crashlooper`

Available tags:
- `latest` - Latest stable release (built from git tags)
- `v*` - Specific version tags (e.g., `v1.0.0`, `v1.0.0-amd64`, `v1.0.0-arm64`)
- `pr-<number>-<sha>` - Pull request builds (e.g., `pr-2-abc1234`)
- `sha-<sha>` - Main branch builds (e.g., `sha-abc1234`)
- All builds include `-amd64` and `-arm64` variants

## Contributing

We use [Conventional Commits](https://www.conventionalcommits.org/) for commit messages and [Semantic Versioning](https://semver.org/) for releases. See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.

### Quick Start for Contributors

```bash
# Check current version
make current-version

# Check what version will be created based on commits
make next-version

# Run tests
make test
```

## Development and Releases

### Automated Versioning and Tagging

This project uses automatic semantic versioning based on commit messages:

- Commits are analyzed using [Conventional Commits](https://www.conventionalcommits.org/)
- When code is merged to `main`, a GitHub Action automatically creates appropriate version tags
- Tags trigger the release workflow to build and publish artifacts

**Commit format determines version bumps:**
- `fix:` commits → PATCH version (v1.0.0 → v1.0.1)
- `feat:` commits → MINOR version (v1.0.0 → v1.1.0)
- `BREAKING CHANGE` or `!` → MAJOR version (v1.0.0 → v2.0.0)

### Automated Docker Publishing

A single unified workflow (`.github/workflows/release.yml`) automatically builds and publishes Docker images in three scenarios:

1. **Pull request builds** (on pull requests):
   - Triggered automatically for all pull requests
   - Images are tagged as `pr-<number>-<short-sha>` (unique per PR commit)
   - Each commit to a PR creates a new tagged image
   - Only Docker images are built (no GitHub releases)

2. **Development builds** (on push to main):
   - Triggered automatically when code is pushed to the `main` branch
   - Images are tagged as `sha-<short-sha>` (unique per commit)
   - Only Docker images are built (no GitHub releases)

3. **Release builds** (on git tags):
   - Triggered when a new version tag is created (e.g., `v1.0.0`)
   - Images are tagged with the version and `latest`
   - Full GitHub release is created with binaries and packages

### Creating a New Release

To publish a new release with Docker images:

```bash
# Create and push a new tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

This will automatically:
- Build binaries for multiple platforms
- Create APK, DEB, and RPM packages
- Build and push multi-arch Docker images (amd64, arm64)
- Create a GitHub release with artifacts
- Update the `latest` Docker tag

### Required GitHub Secrets

The following secrets must be configured in the repository:
- `DOCKER_USERNAME` - Docker Hub username
- `DOCKER_TOKEN` - Docker Hub access token or password
- `GORELEASER_GITHUB_TOKEN` - GitHub token for creating releases
