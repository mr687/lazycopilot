# LazyCopilot

LazyCopilot is a versatile AI-powered tool designed to assist developers with various tasks, starting with generating better commit messages. It leverages the capabilities of GitHub Copilot to provide intelligent suggestions based on your code changes.

## Features

- **Commit Message Generation**: Generates commit messages following the commitizen convention.
- **Authentication**: Seamlessly authenticate with GitHub to access Copilot's capabilities.
- **Versatile Commit Styles**: Choose from different commit title styles such as normal, funny, wise, and trolling.
- **Future Potential**: Plans to expand into other areas of development assistance, making it an indispensable AI tool for developers.

## Installation

### From Source

1. Clone the repository:
    ```sh
    git clone https://github.com/mr687/lazycopilot.git
    cd lazycopilot
    ```

2. Build the project:
    ```sh
    make build
    ```
3. Run
    ```sh
    ./bin/lazycopilot
    ```

### Using Homebrew

1. Add the tap:
    ```sh
    brew tap mr687/lazycopilot
    ```

2. Install LazyCopilot:
    ```sh
    brew install lazycopilot
    ```

### Using Scoop (Windows)

1. Add the bucket:
    ```sh
    scoop bucket add lazycopilot https://github.com/mr687/scoop-lazycopilot
    ```

2. Install LazyCopilot:
    ```sh
    scoop install lazycopilot
    ```

### Binary Release

Binary releases are available for macOS, Linux, and Windows architectures.

1. Download the binary release file from the [releases page](https://github.com/mr687/lazycopilot/releases).
2. Extract the downloaded file.
3. Move the binary to a directory included in your system's PATH.

## Usage

To use LazyCopilot, run the following command in your terminal:

```sh
lazycopilot
```

### Commands

#### `commit`

Generate a commit message using AI.

```sh
lazycopilot commit
```

Options:
- `--stage`: Automatically stage changes if no staged changes are detected.
- `--title-only`: Generate only the commit title without the body.
- `--style`: Specify the style of the commit title. Available options are `normal`, `funny`, `wise`, and `trolling`.
- `--no-commit`: Generate the commit message without committing the changes immediately.

#### `auth`

Authenticate with GitHub to enable Copilot's capabilities.

```sh
lazycopilot auth login
lazycopilot auth logout
```

- `login`: Log in to your GitHub account.
- `logout`: Log out from your GitHub account.

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
