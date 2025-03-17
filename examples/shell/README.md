# Shell key bindings

These widgets bind <kbd>Ctrl</kbd>+<kbd>S</kbd> to `snipkit pick`, so you can
fuzzy-pick a snippet and drop it straight onto the command line without leaving
what you were typing.

| Shell | File | How to load |
| ----- | ---- | ----------- |
| bash  | `bash-widget.sh`   | `source .../bash-widget.sh` in `~/.bashrc` |
| zsh   | `zsh-widget.zsh`   | `source .../zsh-widget.zsh` in `~/.zshrc`  |
| fish  | `fish-widget.fish` | `source .../fish-widget.fish` in `config.fish` |

`snipkit pick` uses [`fzf`](https://github.com/junegunn/fzf) when it is on your
`PATH`, and falls back to a numbered prompt otherwise.
