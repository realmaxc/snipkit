# snipkit bash widget — bind Ctrl-S to pick a snippet and insert it at the
# cursor. Add to ~/.bashrc:
#
#   source /path/to/snipkit/examples/shell/bash-widget.sh
#
_snipkit_pick() {
    local selection
    selection="$(snipkit pick 2>/dev/tty)" || return
    [ -z "$selection" ] && return
    READLINE_LINE="${READLINE_LINE:0:$READLINE_POINT}${selection}${READLINE_LINE:$READLINE_POINT}"
    READLINE_POINT=$(( READLINE_POINT + ${#selection} ))
}
bind -x '"\C-s": _snipkit_pick' 2>/dev/null
