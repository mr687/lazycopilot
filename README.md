# LazyCopilot

LazyCopilot is a tool to help you write better commit messages by leveraging AI capabilities. It integrates with GitHub Copilot to generate commit messages based on your code changes.

## Features

- Generates commit messages following the commitizen convention.
- Uses GitHub Copilot for generating commit messages.
- Easy to use CLI interface.

## Requirements

- NVIM (Neovim) installed.
- The `copilot.lua` plugin installed and configured to sign in to GitHub.

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