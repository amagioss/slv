---
sidebar_position: 10
---

# Contributing to SLV

Thank you for your interest in contributing to SLV! This guide will help you get started with contributing to the project.

---

## How to Contribute

There are many ways to contribute to SLV:

- **🐛 Report bugs** - Open an issue describing the problem
- **💡 Suggest features** - Share your ideas for improvements
- **📝 Improve documentation** - Help make the docs better
- **💻 Write code** - Fix bugs or add new features
- **🧪 Write tests** - Improve test coverage
- **🔍 Review pull requests** - Help review and test changes

---

## Getting Started

### Prerequisites

Before you begin, ensure you have the following installed:

- **Go** 1.24 or later ([Download](https://go.dev/dl/))
- **Git** ([Download](https://git-scm.com/downloads))
- **Node.js** 18+ and npm (for documentation contributions)

### Development Setup

1. **Fork the repository**
   ```bash
   # Fork the repository on GitHub, then clone your fork
   git clone https://github.com/YOUR_USERNAME/slv.git
   cd slv
   ```

2. **Add the upstream remote**
   ```bash
   git remote add upstream https://github.com/amagioss/slv.git
   ```

3. **Install dependencies**
   ```bash
   # Go dependencies will be downloaded automatically
   go mod download
   
   # For documentation development
   cd website
   npm install
   ```

4. **Build the project**
   ```bash
   # Build the main CLI
   go build -o slv ./internal/app/main.go
   
   # Or install it
   go install ./internal/app/main.go
   ```

---

## Development Workflow

### 1. Create a Branch

Always create a new branch for your changes:

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-bug-fix
```

**Branch naming conventions:**
- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation changes
- `refactor/` - Code refactoring
- `test/` - Test additions or improvements

### 2. Make Your Changes

- Write clean, readable code
- Follow Go conventions and best practices
- Add comments for complex logic
- Update tests if needed
- Update documentation if your changes affect user-facing features

### 3. Test Your Changes

```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./internal/core/vaults

# Run with verbose output
go test ./... -v
```

### 4. Commit Your Changes

Write clear, descriptive commit messages:

```bash
git add .
git commit -m "Add feature: description of what you added"
```

**Commit message guidelines:**
- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters
- Reference issues and pull requests when applicable

**Example:**
```
Add support for Azure Key Vault integration

- Implement Azure KMS provider
- Add authentication flow
- Update documentation

Fixes #123
```

### 5. Push and Create a Pull Request

```bash
git push origin feature/your-feature-name
```

Then open a Pull Request on GitHub with:
- A clear title and description
- Reference to related issues
- Screenshots (if applicable)
- Testing instructions

---

## Code Contributions

### Project Structure

```
slv/
├── internal/
│   └── app/main.go   # Entry Point
│   ├── cli/          # CLI commands
│   ├── core/         # Core functionality
│   ├── tui/          # Terminal UI
│   └── k8s/          # Kubernetes operator
├── action/           # GitHub Actions
├── website/          # Documentation site
└── slv.go           # Library Functions
```

### Code Style

- Follow [Effective Go](https://go.dev/doc/effective_go) guidelines
- Use `gofmt` to format your code
- Run `go vet` to check for common mistakes
- Use meaningful variable and function names
- Keep functions small and focused

---

## Documentation Contributions

The documentation is built with [Docusaurus](https://docusaurus.io/).

### Setup

```bash
cd website
npm install
```

### Development

```bash
# Start the development server
npm start

# Build the site
npm run build

# Clear cache
npm run clear
```

### Documentation Structure

- **Markdown files** in `website/docs/` - Main documentation content
- **Components** in `website/src/components/` - React components
- **Styling** in `website/src/css/` - Custom CSS

### Writing Documentation

- Use clear, concise language
- Include code examples where helpful
- Add cross-references to related pages
- Keep examples up-to-date with the codebase
- Test all code examples

### Documentation Style Guide

- Use active voice
- Write for both beginners and advanced users
- Include "See Also" sections for related topics
- Add prerequisites when needed
- Use consistent formatting

---

## Testing

### Writing Tests

- Write tests for new features
- Update tests when fixing bugs
- Aim for good test coverage
- Use table-driven tests when appropriate

### Test Structure

```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {
            name:     "valid input",
            input:    "test",
            expected: "result",
            wantErr:  false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run tests for a specific package
go test ./internal/core/vaults
```

---

## Pull Request Process

1. **Update your branch**
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Ensure tests pass**
   ```bash
   go test ./...
   ```

3. **Update documentation** if your changes affect user-facing features

4. **Create the PR** with:
   - Clear title and description
   - Reference to related issues
   - Testing instructions
   - Screenshots (if applicable)

5. **Respond to feedback** and make requested changes

6. **Keep your PR focused** - One feature or fix per PR

---

## Code Review Guidelines

### For Contributors

- Be open to feedback
- Respond to comments promptly
- Make requested changes clearly
- Ask questions if something is unclear

### For Reviewers

- Be constructive and respectful
- Explain the reasoning behind suggestions
- Approve when changes look good
- Test the changes when possible

---

## Reporting Bugs

When reporting bugs, please include:

- **Description** - What happened?
- **Steps to reproduce** - How can we reproduce it?
- **Expected behavior** - What should have happened?
- **Actual behavior** - What actually happened?
- **Environment** - OS, Go version, SLV version
- **Logs** - Relevant error messages or logs

---

## Suggesting Features

When suggesting features, please include:

- **Use case** - What problem does this solve?
- **Proposed solution** - How should it work?
- **Alternatives** - Other solutions you've considered
- **Additional context** - Any other relevant information

---

## Questions?

- **GitHub Discussions** - For questions and discussions
- **GitHub Issues** - For bug reports and feature requests
- **Documentation** - Check the docs at [slv.sh](https://slv.sh)

---

## License

By contributing to SLV, you agree that your contributions will be licensed under the MIT License.

---

## See Also

- [Quick Start Guide](/docs/quick-start) - Get started with SLV
- [Command Reference](/docs/command-reference/vault/get) - Learn about SLV commands
- [Overview](/docs/overview) - Learn about SLV

