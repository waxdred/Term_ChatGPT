![GitHub Workflow Status](https://github.com/waxdred/Term_ChatGPT/actions/workflows/default.yml/badge.svg)
![Go](https://img.shields.io/badge/Made%20with%20Goland-blue.svg?style=for-the-badge&logo=goland)
# Neovim_ChatGPT

`Plugin working but still in Progress`

`ChatGPT` is a Neovim plugin or app for your terminal that allows you to interact with OpenAI's GPT-3 language model.
With `ChatGPT`, you can ask questions and get answers from GPT-3 in real-time.

`Plugin create in Goland using charmbracelet/bubbles for the view`

## Installation

- Set environment variable called `$OPENAI_API_KEY` which you can [obtain here](https://beta.openai.com/account/api-keys).

- Neovim Plugin require
```
require github.com/voldikss/vim-floaterm
```

```lua
-- Packer
use { "github.com/waxdred/Term_ChatGPT" }
vim-floaterm is for run chatGpt on your terminal
use { "github.com/voldikss/vim-floaterm" }
```

## Configuration


## Usage

Plugin exposes following commands on Neovim:
- `:ChatGPT` command which opens interactive window.

Available keybindings are:
- `<Esc>` to close chat window.
- `scroll mouse` scroll up chat window.
- `scroll mouse` scroll down chat window.
- `<C-y>` to copy/yank last answer.
- `<C-n>` Start new session in Progress.
- `<Tab>` Cycle over windows.
- Setting Window
- `<C-k>` [Change value selection] up the value
- `<C-j>` [Change value selection] down the value
- `<up>` [Navigate] with arrow
- `<down>` [Navigate] with arrow
