# LazyCopilot

LazyCopilot is a versatile AI-powered tool designed to assist developers with various tasks, starting with generating better commit messages. It leverages the capabilities of GitHub Copilot to provide intelligent suggestions based on your code changes.

## Features

- **Commit Message Generation**: Generates commit messages following the commitizen convention.
- **Authentication**: Seamlessly authenticate with GitHub to access Copilot's capabilities.
- **Versatile Commit Styles**: Choose from different commit title styles such as normal, funny, wise, and trolling.
- **Future Potential**: Plans to expand into other areas of development assistance, making it an indispensable AI tool for developers.

## Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/mr687/lazycopilot.git
    cd lazycopilot
    ```

2. Build the project:
    ```sh
    make build
    ```

## Usage

To use LazyCopilot, run the following command in your terminal:

```sh
bin/lazycopilot --path /path/to/your/git/repository
```

### Commands

#### `commit`

Generate a commit message using AI.

```sh
bin/lazycopilot commit --path /path/to/your/git/repository
```

Options:
- `--stage`: Stage changes if no staged changes are detected.
- `--title-only`: Generate only the commit title.
- `--style`: Style of the commit title: `normal`, `funny`, `wise`, `trolling`.
- `--no-commit`: Do not commit the generated content immediately.

#### `auth`

Authenticate with GitHub.

```sh
bin/lazycopilot auth login
bin/lazycopilot auth logout
```

## Future Plans

LazyCopilot aims to evolve into a comprehensive AI assistant for developers. Future features may include:

- **Code Review Assistance**: Provide AI-driven code review suggestions.
- **Documentation Generation**: Automatically generate documentation based on code comments and structure.
- **Bug Detection**: Identify potential bugs and suggest fixes.
- **Code Refactoring**: Offer intelligent code refactoring suggestions.

## Screenshots

Here are some screenshots demonstrating LazyCopilot in action:

### Commit Message Generation

![Commit Message Generation](screenshots/commit-message-generation.png)

## Contributing

Contributions are welcome! Please read the [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to contribute to this project.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
