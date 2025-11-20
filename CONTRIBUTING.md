# Contributing to Quazaar

Thank you for your interest in contributing to Quazaar! We welcome contributions from the community. This document provides guidelines and information for contributors.

## 🚀 Ways to Contribute

### Report Issues

- Use GitHub Issues to report bugs, request features, or ask questions
- Check existing issues first to avoid duplicates
- Provide detailed information including:
  - Steps to reproduce
  - Expected vs actual behavior
  - Environment details (OS, Go version, etc.)
  - Error messages/logs

### Submit Pull Requests

- Fork the repository
- Create a feature branch from `main`
- Make your changes
- Test thoroughly
- Submit a pull request with a clear description

### Improve Documentation

- Fix typos or unclear explanations
- Add missing documentation
- Create tutorials or examples
- Update API documentation

## 🛠️ Development Setup

### Prerequisites

- Go 1.21 or higher
- Git
- (Optional) `playerctl` for music integration testing

### Setup Steps

1. **Clone and setup**

   ```bash
   git clone https://github.com/codershubinc/Quazaar.git
   cd Quazaar
   go mod download
   ```

2. **Build the project**

   ```bash
   go build -o quazaar ./cmd/server
   ```

3. **Run tests**

   ```bash
   go test ./...
   ```

4. **Run the server**

   ```bash
   ./quazaar
   ```

## 📝 Code Style Guidelines

### Go Code

- Follow standard Go formatting (`go fmt`)
- Use `gofmt` and `goimports` for consistent formatting
- Follow Go naming conventions (camelCase for private, PascalCase for exported)
- Write clear, concise comments for exported functions/types
- Use meaningful variable and function names
- Keep functions small and focused on single responsibilities

### Commit Messages

- Use clear, descriptive commit messages
- Start with a verb in imperative mood (e.g., "Add", "Fix", "Update")
- Keep first line under 50 characters
- Add detailed description if needed

Examples:

```text
Add Spotify artist WebSocket handlers

- Implement getArtistInfo and followArtist WebSocket endpoints
- Add proper error handling for API failures
- Update frontend with artist controls
```

### Code Structure

- Keep packages focused and modular
- Use `internal/` for private packages
- Follow the established project structure
- Add tests for new functionality

## 🧪 Testing

### Unit Tests

- Write tests for new functions and methods
- Use Go's built-in testing framework
- Place test files alongside source files (e.g., `handler_test.go`)
- Aim for good test coverage

### Integration Tests

- Test WebSocket connections and message handling
- Test API endpoints
- Use the provided test scripts in `tests/`

### Manual Testing

- Test with the web client in `temp/web/`
- Verify functionality across different browsers
- Test on different operating systems if possible

## 🔄 Pull Request Process

1. **Create a branch** from `main`:

   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following the guidelines above

3. **Test thoroughly**:

   - Run existing tests: `go test ./...`
   - Build successfully: `go build ./cmd/server`
   - Test manually with the web client

4. **Update documentation** if needed:

   - Update README if adding new features
   - Add API documentation for new endpoints
   - Update changelog

5. **Commit your changes**:

   ```bash
   git add .
   git commit -m "Add your descriptive commit message"
   ```

6. **Push and create PR**:

   ```bash
   git push origin feature/your-feature-name
   ```

   Then create a pull request on GitHub

7. **Address feedback** from reviewers and make necessary changes

## 📋 Pull Request Checklist

Before submitting a PR, ensure:

- [ ] Code builds successfully (`go build ./cmd/server`)
- [ ] All tests pass (`go test ./...`)
- [ ] Code follows style guidelines (`go fmt`, `go vet`)
- [ ] No linting errors (`golangci-lint run` if available)
- [ ] Documentation updated if needed
- [ ] Commit messages are clear and descriptive
- [ ] Changes tested manually
- [ ] No sensitive information committed

## 🎯 Areas for Contribution

### High Priority

- Bug fixes and stability improvements
- Security enhancements
- Performance optimizations
- Better error handling

### Medium Priority

- New features and integrations
- UI/UX improvements
- Documentation improvements
- Test coverage expansion

### Low Priority

- Code refactoring and cleanup
- New themes and customization options
- Additional platform support

## 📞 Getting Help

- **Issues**: Use GitHub Issues for bugs and feature requests
- **Discussions**: Use GitHub Discussions for questions and general discussion
- **Documentation**: Check the `docs/` folder for detailed guides

## 📜 Code of Conduct

Please be respectful and constructive in all interactions. We follow a standard code of conduct to ensure a welcoming environment for all contributors.

## 🙏 Recognition

Contributors will be acknowledged in the project documentation and release notes. Thank you for helping make Quazaar better!

---

Made with ❤️ by the Quazaar community
