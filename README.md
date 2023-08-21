# spux

A cli tool to make creating tmux development environments dead easy.

## The problem

I use neovim as my primary development tool, and tmux to organize my
terminal while working projects. In general, this setup is awesome. 
However, every time I reboot my computer, I lose my tmux sessions.

I started to write bash scripts to recreate and 'save' my project 
workspaces. But, I quickly discovered that writing these scripts is a hassle.
Working with the tmux cli syntax takes my focus off of actually designing 
my dev environment, and tmux's target system can make these scripts hard to
maintain. This is especially painful when you just want to set up (and save)
an environment in the early stages of prototyping.

## The solution

```yaml
space: "spux"
root: "/home/sawyer/dev/personal/spux"
windows:
  - name: "shell"
  - name: "vi"
    init: "nvim . {enter}"
    cmds:
      - cmd: "{clear} go build spux.go {enter}"
        tgt: "shell"
```
