# Contributing to Quazaar

We love your input! We want to make contributing to Quazaar as easy and transparent as possible, whether it's:

- Reporting a bug
- Discussing the current state of the code
- Submitting a fix
- Proposing new features
- Becoming a maintainer

## Development Process

We use GitHub to host code, to track issues and feature requests, and to accept pull requests.

1. Fork the repo and create your branch from `main`.
2. If you've added code that should be tested, add tests.
3. If you've changed APIs, update the documentation.
4. Ensure the test suite passes.
5. Make sure your code lints.
6. Issue that pull request!

## Development Setup

### Prerequisites

- **Go**: Version 1.21 or higher
- **Git**: For version control
- **Playerctl**: (Optional) For testing media controls on Linux

### Getting Started

1. **Clone the repository**

   ```bash
   git clone https://github.com/codershubinc/Quazaar.git
   cd Quazaar
   ```

2. **Install dependencies**

   ```bash
   go mod download
   ```

3. **Run the server**

   ```bash
   go run cmd/server/main.go
   ```

4. **Run tests**

   ```bash
   go test ./...
   ```

## Pull Request Process

1. Update the `README.md` with details of changes to the interface, this includes new environment variables, exposed ports, useful file locations and container parameters.
2. Increase the version numbers in any examples files and the README.md to the new version that this Pull Request would represent. The versioning scheme we use is [SemVer](http://semver.org/).
3. You may merge the Pull Request in once you have the sign-off of two other developers, or if you do not have permission to do that, you may request the second reviewer to merge it for you.

## Code Style

- **Formatting**: We use `gofmt` to format our code. Please run it before committing.
- **Linting**: We use `golangci-lint`. You can run it locally with `golangci-lint run`.
- **Comments**: Document exported functions and complex logic.

## License

By contributing, you agree that your contributions will be licensed under its MIT License.
