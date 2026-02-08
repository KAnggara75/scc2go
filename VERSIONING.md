# Semantic Versioning Guide

This project follows [Semantic Versioning 2.0.0](https://semver.org/) with automated version bumping based on [Conventional Commits](https://www.conventionalcommits.org/).

## Version Format

Versions follow the format: `vMAJOR.MINOR.PATCH`

- **MAJOR**: Breaking changes (incompatible API changes)
- **MINOR**: New features (backwards-compatible)
- **PATCH**: Bug fixes (backwards-compatible)

## Commit Message Format

Use conventional commit messages to automatically determine version bumps:

### Major Version Bump (Breaking Changes)

```
feat!: remove deprecated GetConfig function

BREAKING CHANGE: GetConfig has been removed, use GetEnv instead
```

Or:

```
refactor!: change function signature for GetEnv
```

### Minor Version Bump (New Features)

```
feat: add support for custom timeout configuration
```

```
feat(cache): implement Redis caching support
```

### Patch Version Bump (Bug Fixes)

```
fix: correct error handling in getSCC function
```

```
fix(auth): resolve token validation issue
```

```
perf: improve config parsing performance
```

```
refactor: simplify setIfNotExists logic
```

### Other Commit Types (Also Patch Bump)

Any commit without a conventional prefix will trigger a patch bump:

```
update documentation
improve test coverage
```

## How It Works

1. **On Push to Main/Master**: The GitHub Actions workflow analyzes commits since the last tag
2. **Version Calculation**: Based on commit messages, it determines the appropriate version bump
3. **Tag Creation**: Creates and pushes a new git tag with the calculated version
4. **GitHub Release**: Automatically creates a GitHub release with changelog

## Examples

### Current version: v0.1.6

| Commit Message | New Version | Reason |
|---------------|-------------|---------|
| `feat: add new feature` | v0.2.0 | Minor bump (new feature) |
| `fix: bug fix` | v0.1.7 | Patch bump (bug fix) |
| `feat!: breaking change` | v1.0.0 | Major bump (breaking) |
| `docs: update README` | v0.1.7 | Patch bump (default) |

## Manual Tagging

If you need to create a tag manually:

```bash
# Create an annotated tag
git tag -a v1.0.0 -m "Release version 1.0.0"

# Push the tag
git push origin v1.0.0
```

## Best Practices

1. **Use Conventional Commits**: Always prefix commits with type (feat, fix, etc.)
2. **Be Specific**: Use scopes to indicate what changed: `feat(cache): ...`
3. **Breaking Changes**: Always use `!` or `BREAKING CHANGE:` for incompatible changes
4. **Descriptive Messages**: Write clear commit messages that explain the change

## CI/CD Integration

The workflow runs on:
- ✅ Push to main/master branches
- ✅ Pull requests (build and test only, no tagging)

### Workflow Steps

1. **Build**: Compiles the code with multiple Go versions
2. **Test**: Runs tests with coverage reporting
3. **Lint**: Runs pre-commit hooks and linters
4. **Version**: Determines next version from commits
5. **Tag**: Creates and pushes the new version tag
6. **Release**: Creates GitHub release with changelog

## Viewing Version History

```bash
# List all tags
git tag -l

# Show latest tag
git describe --tags --abbrev=0

# View changelog between versions
git log v0.1.5..v0.1.6 --oneline
```

## Migration from Old System

The old system used date-based versioning: `v0.yy.m-dhhmm`

The new system uses semantic versioning: `vMAJOR.MINOR.PATCH`

Old tags are preserved for historical reference.
