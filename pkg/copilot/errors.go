package copilot

import "errors"

var ErrGithubTokenNotSet = errors.New("no GitHub token found, please use `:Copilot auth` to set it up from copilot.lua or `:Copilot setup` for copilot.vim")
