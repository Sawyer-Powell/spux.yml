# spux

A cli tool to make creating and managing tmux development environments 
dead simple.

### Table of Contents

- [The Problem](#the-problem)
- [The Solution](#the-solution)
- [Installation](#installation)
- [The spux.yml file](#the-spux.yml-file)
- [Making use of SPUX_CMD environment variables](#making-use-of-spux_cmd-environment-variables)

## The Problem

Doing development in vim and tmux is awesome, but having to recreate your
carefully crafted dev environment every time you reboot is not.

I started to write bash scripts to recreate and 'save' my project 
workspaces. But, writing these scripts can be frustrating.
I found that scripting with the syntax of the tmux cli can take my focus off of 
*actually designing my dev environment*, and tmux's target system can make 
these scripts hard to maintain. This is especially painful when you just want 
to set up (and save) an environment in the early stages of prototyping.

To me, having as low a barrier as possible from idea, to implementing it in a 
quality development environment, is important.

## The Solution

`spux.yml`
```yaml
space: "spux"
root: "/home/sawyer/dev/personal/spux"
windows:
  - name: "shell"
    panes:
        - name: "git"
          init: "lazygit {enter}"
        - name: "zsh"
    layout: "even-horizontal"
  - name: "vi"
    init: "nvim . {enter}"
    cmds:
      - cmd: "{clear} go build spux.go {enter}"
        tgt: "zsh"
```

By creating a `spux.yml` file in your project's root directoy, you can just run
```bash
spux
```
and spux will:
1. Take your `yml` file and generate the bash script below
2. Save that bash script into `~/.config/spux/bin`
3. Run the bash script and attach your terminal to the newly created session

```bash
#!/bin/bash
cd /home/sawyer/dev/personal/spux
tmux new-session -d -s spux
tmux rename-window -t spux:0 shell
tmux send-keys -t spux:0.0 'lazygit' C-m
tmux split-window -t spux:0.0
tmux send-keys -t spux:0.1
tmux select-layout -t spux:0 even-horizontal
tmux setenv -t spux:1 SPUX_CMD_0 "tmux send-keys -t spux:0.1 C-l 'go build spux.go' C-m ;tmux select-window -t spux:0;tmux select-pane -t spux:0.1"
tmux new-window -t spux:1 -n vi
tmux send-keys -t spux:1 'nvim .' C-m
```

## Installation

### Requirements

Spux is only compatible with linux/unix machines which have tmux installed

Spux is not available for windows machines.

### Building and installing

Copy and use the `spux` binary provided in the root of this package at
your own risk.

I recommend building it yourself:

1. Make sure you have golang installed, if not install [here](https://go.dev/doc/install)
2. Run the script below

```bash
git clone https://github.com/Sawyer-Powell/spux.yml
cd spux.yml
go get -u ./... # grabs all the dependencies for this project
go build spux.go
sudo cp ./spux /usr/bin # or somewhere else on your path, I keep my spux binary in ~/.local/bin
```

## The spux.yml file 

Here is a full specification of the spux.yml file:

```yaml
space: "the name of the tmux session" # required
root: "the root directory of this session" # required
windows:
    - name: "the name of the window (also functions as its id for spux's target system)" # required to define a window
      init: "{clear} script you want this window to run when its created {enter}" # optional
      # {clear} clears the terminal
      # {enter} runs the command preceding it in the terminal
      # {interrupt} sends Ctrl-C to the terminal
    - name: "my other window"
      panes:
        - name: "the id of this pane (for spux's targetting system)" # required to define a pane
          init: "{clear} launch vim {enter}" # optional
          cmds: 
              # The below will generate a script which will run against the
              # 'frontend_server' pane. spux will bind this script to the $SPUX_CMD_0
              # environment variable, which exported to this pane.
              # While in this pane, simply execute the command in 
              # $SPUX_CMD_0 and it will run in the target specified in "tgt", and
              # will automatically switch your tmux focus.
            - cmd: "{interrupt} {clear} launch backend server {enter}"
              tgt: "frontend_server"
              # This will create a script and assign it to $SPUX_CMD_1 in this pane
            - cmd: "another command"
              tgt: "frontend_server"
        - name: "frontend_server"
          init: "{clear} launch backend server {enter}"
      layout: "even-horizontal | even-vertical | main-horizontal | main-vertical | tiled" #optional, specifies the tmux layout for this windows panes
```

## Making use of SPUX_CMD environment variables

Here's what's in my neovim config

```lua
vim.keymap.set('n', '<leader><enter>', function()
	local status, output, exit_code = os.execute(vim.env.SPUX_CMD_0)
	if status then
		print("$SPUX_CMD_0 executed successfully.")
	else
		print("Error executing $SPUX_CMD_0.")
	end
end)
```

So, referencing the script provided [here](#the-solution),
whenever I press `space` (my leader key) then `enter` in my neovim open in
the `vi` window, it runs `go build spux.go` in my `zsh` pane.
