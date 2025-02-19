# Contributing to convo-ai-go-server

First off, thank you for considering contributing to **convo-ai-go-server**! 🎉 Your contributions are greatly appreciated.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How to Contribute](#how-to-contribute)
  - [Reporting Bugs](#reporting-bugs)
  - [Suggesting Enhancements](#suggesting-enhancements)
  - [Submitting Pull Requests](#submitting-pull-requests)
- [Style Guides](#style-guides)
  - [Coding Style](#coding-style)
  - [Commit Messages](#commit-messages)
- [Testing](#testing)
- [License](#license)

## Code of Conduct

Please read and follow the [Code of Conduct](CODE_OF_CONDUCT.md) to ensure a welcoming and respectful environment for all contributors.

## How to Contribute

### Reporting Bugs

If you find a bug in the project, please open an issue in the [GitHub Issues](https://github.com/yourusername/convo-ai-go-server/issues) section. Make sure to include:

- A clear and descriptive title
- A step-by-step description of how to reproduce the bug
- Expected and actual behavior
- Any relevant screenshots or logs

### Suggesting Enhancements

Have an idea for improving the project? We'd love to hear it! Please open an issue with the tag `enhancement` and include:

- A clear and descriptive title
- A detailed description of the improvement
- Any relevant examples or use cases

### Submitting Pull Requests

Contributions are welcome! To submit a pull request, follow these steps:

1. **Fork the Repository**

   - Click the "Fork" button at the top of the repository page.

2. **Clone Your Fork**

   ```bash
   git clone https://github.com/yourusername/convo-ai-go-server.git
   cd convo-ai-go-server
   ```

3. **Create a New Branch**

   ```bash
   git checkout -b feature/your-feature-name
   ```

4. **Make Your Changes**

   - Implement your feature or bug fix.
   - Ensure your code adheres to the project's coding style.

5. **Run Tests**

   ```bash
   go test ./...
   ```

6. **Commit Your Changes**

   ```bash
   git commit -m "Add feature: your feature description"
   ```

7. **Push to Your Fork**

   ```bash
   git push origin feature/your-feature-name
   ```

8. **Open a Pull Request**
   - Navigate to the original repository and click "New Pull Request".
   - Provide a clear title and description of your changes.

## Style Guides

### Coding Style

- Follow [Go's official style guidelines](https://golang.org/doc/effective_go.html).
- Use `gofmt` to format your code:
  ```bash
  gofmt -w .
  ```

### Commit Messages

- Use the present tense ("Add feature" not "Added feature").
- Use the imperative mood ("Fix bug" not "Fixes bug").
- Limit the first line to 50 characters.
- Provide a detailed description if necessary.

## Testing

Ensure all tests pass before submitting a pull request.

```bash
go test ./...
```

Write tests for new features and bug fixes to maintain the project's reliability.

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).

For any questions or further assistance, feel free to reach out via the [issues page](https://github.com/yourusername/convo-ai-go-server/issues).

Thank you for your contributions! 🙏
