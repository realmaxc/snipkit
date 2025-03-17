# snipkit zsh widget — bind Ctrl-S to pick a snippet and insert it at the
# cursor. Add to ~/.zshrc:
#
#   source /path/to/snipkit/examples/shell/zsh-widget.zsh
#
snipkit-pick-widget() {
    local selection
    selection="$(snipkit pick 2>/dev/tty)" || return
    [[ -z "$selection" ]] && return
    LBUFFER="${LBUFFER}${selection}"
    zle reset-prompt
}
zle -N snipkit-pick-widget
bindkey '^s' snipkit-pick-widget
