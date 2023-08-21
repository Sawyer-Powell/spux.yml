# spux

A cli tool to make creating and managing tmux development environments 
dead simple.

## The Problem

Doing development in vim and tmux is awesome, but having to recreate your
project setup every time you reboot is not.

I started to write bash scripts to recreate and 'save' my project 
workspaces. But, writing these scripts can be frustrating.
I found that the syntax of scripting with the tmux cli can take my focus off of 
*actually designing my dev environment*, and tmux's target system can make 
these scripts hard to maintain. This is especially painful when you just want 
to set up (and save) an environment in the early stages of prototyping.

To me, having as low a barrier as possible from idea, to implementing it in a 
quality development environment, is important.

## The Solution

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
3. Run the bash script and attach your terminal to newly created session

```bash
#!/bin/bash
cd /home/sawyer/dev/personal/spux
tmux new-session -d -s spux
tmux rename-window -t spux:0 shell
tmux send-keys -t spux:0.0 C-l 'lazygit' C-m
tmux split-window -t spux:0.0
tmux send-keys -t spux:0.1 C-l
tmux select-layout -t spux:0 even-horizontal
tmux setenv -t spux:1 SPUX_CMD_0 "tmux send-keys -t spux:0.1 C-l 'go build spux.go' C-m ;tmux select-window -t spux:0;tmux select-pane -t spux:0.1"
tmux new-window -t spux:1 -n vi
tmux send-keys -t spux:1 C-l 'nvim .' C-m
```
