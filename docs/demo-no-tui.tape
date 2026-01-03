Output demo-no-tui.gif
Set Shell fish
Set Theme "Catppuccin Mocha"
Set Width 1920
Set Height 1080

Sleep 1s
Type "./bin/optinix get --no-tui | grep '^Option [0-9]' | cut -d' ' -f3 | fzf --preview='./bin/optinix get --no-tui \'{}\' '"
Enter
Sleep 3s

Type "podman"
Sleep 3s
Enter
Sleep 3s

