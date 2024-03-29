![GitHub Workflow Status](https://github.com/waxdred/Term_ChatGPT/actions/workflows/default.yml/badge.svg)
![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)
# Neovim_ChatGPT

`ChatGPT` is a Neovim plugin or app for your terminal that allows you to interact with OpenAI's GPT-3 language model.
With `ChatGPT`, you can ask questions and get answers from GPT-3 in real-time.
and save your sessions for future use


The Neovim_ChatGPT plugin is a powerful tool that allows users to interact with OpenAI's GPT-3 language model directly from their Neovim editor or terminal. With this plugin, users can ask any question and get accurate, relevant answers in real-time. The plugin also allows users to save their sessions for future use, making it easy to continue conversations or pick up where they left off.

The main benefits of this plugin include:

- Convenient and easy access to the power of GPT-3 language model from within Neovim or terminal.
- Real-time answers to any question, making it a valuable tool for research, writing, or any other task that        requires quick and accurate information.
- The ability to save and continue sessions for later use, making it a more efficient tool for tasks that require multiple interactions with the language model.
- A simple and user-friendly interface that makes it easy for anyone to use.
- Plugin create in Goland using library `charmbracelet/bubbles`

## Table of Contents
- [Requirements](#Requirements)
- [Installation](#Installation)
- [Usage](#Usage)
- [Setting](#Setting)
## App chatGpt can be use on your terminal or on Neovim </br></br>

![](https://imgur.com/KCCKaR6.png)

## Requirements
- Neovim: The plugin is designed to work with Neovim, so you'll need to have Neovim installed on your machine. You should have at least version 0.5.0 of Neovim, but it's recommended to use the latest version for the best experience.

- Golang: The plugin is built using the Golang programming language, so you'll need to have Golang installed on your machine. You should have at least version 1.15 of Golang, but it's recommended to use the latest version for the best experience.

- OpenAI API Key: To use the plugin, you'll need to have an API key from OpenAI, which you can obtain by signing up for an account on the OpenAI website. You will need to set an environment variable called $OPENAI_API_KEY with the value of the API key.

- Vim-floaterm: You also need to require the vim-floaterm library, this is for run chatGpt on your terminal

- Packer: You need to use the package manager for installing the plugin, such as vim-plug, dein.vim or Vundle.vim

![](https://i.imgur.com/YIyqUFL.gif)

## Installation

- Set environment variable called `$OPENAI_API_KEY` which you can [obtain here](https://beta.openai.com/account/api-keys).

- Neovim Plugin require
```
require github.com/voldikss/vim-floaterm
```

```lua
-- Packer
use { "waxdred/Term_ChatGPT" }

-- After install with packer use:
:ChatGPTInstall

-- vim-floaterm is for run chatGpt on your terminal
use { "voldikss/vim-floaterm" }
```

- Using on your terminal
```
git clone https://github.com/waxdred/Term_ChatGPT.git 
cd Term_chatGPT
./install.sh
cd bin
./chatGPT
```

## Usage

Plugin exposes following commands on Neovim:
- `:ChatGPT` command which opens interactive window.
- on terminal can use the flag:
```
  flag list:
    flag      | default
    -temp     | 0
    -top      | 0.5
    -freg     | 0.5
    -pres     | 0.5
    -token    | 100
    
  exemple:
    ./chatGPT -temp=0 -top=1 -freg=0.5 -pres=0.6 -token=30
```
## Setting
#### Setting viewPort
- The setting can be change easily </br> </br>
![](https://i.imgur.com/5TFJflJ.gif) </br></br>
- temperature:`0 | 1`
- topP: `0 | 1`
- frequency: `-2 | 2`
- presence: `-2 | 2`
- maxToken: `0 | 4000`
- for more informations about it [openAI](https://beta.openai.com/docs/guides/completion/prompt-design)
<details>
<summary>Available keybindings</summary>
<br>
  <lu>
    <li>< Esc > to close chat window.</li>
    <li>scroll mouse scroll up chat window.</li>
    <li>scroll mouse scroll down chat window.</li>
    <li>< C-y > to copy/yank last answer.</li>
    <li>< C-n > Start new session.</li>
    <li>< Tab > Cycle over windows.</li>
    <li>Setting Window</li>
    <li>< + > [Change value selection] up the value</li>
    <li>< - > [Change value selection] down the value</li>
    <li>< C-k > [Navigate] with arrow</li>
    <li>< C-j > [Navigate] with arrow</li>
    <li>Session Window</li>
    <li>< Enter > Select Session</li>
    <li>< C-d > deleted Session</li>
    <li>< C-r> rename Session</li>
  </lu>
</details>

## Contributing
Contributions to this project are welcome. If you'd like to contribute, please fork the repository and make your changes. Then, open a pull request and I'll review your changes.

## License
This project is licensed under the MIT License.
