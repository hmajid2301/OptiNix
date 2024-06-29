Output demo-no-tui.gif
Set Shell fish
Set Theme "Catppuccin Mocha"
Set Width 1920
Set Height 1080

Sleep 1s
Type "./bin/optinix --no-tui | rg 'Name: '| cut -d' ' -f2 | fzf --preview='./bin/optinix --no-tui \'{}\' '"
Enter
Sleep 3s

Type "podman"
Sleep 3s
Enter
Sleep 3s

