# snipkit fish widget — bind Ctrl-S to pick a snippet and insert it at the
# cursor. Add to ~/.config/fish/config.fish:
#
#   source /path/to/snipkit/examples/shell/fish-widget.fish
#
function _snipkit_pick
    set -l selection (snipkit pick 2>/dev/tty)
    test -z "$selection"; and return
    commandline -i -- $selection
end
bind \cs _snipkit_pick
