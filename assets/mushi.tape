Output mushi.tape.gif
Set Shell "zsh"
Set Theme "TokyoNightStorm"
Set FontFamily "CaskaydiaCove Nerd Font"
Set FontSize 30
Set Width 1500
Set Height 900
Set Padding 10
Set TypingSpeed 0.2

Hide
Type@0ms "mkdir -p mushi_test"
Enter
Type@0ms "cd mushi_test"
Enter

Sleep 2.0
Show
Type "ls -a"
Enter
Type "mushi create -i"
Sleep 0.2
Enter
Type "/go"
Sleep 1.0
Enter
Sleep 1.0
Type "ls -a"
Enter
Sleep 0.5
Type "cat .gitignore"
Enter
Sleep 2.0
Hide
Type "cd .."
Enter
Type "rm -rf mushi_test"
Enter
Type@0ms ":qa!" Enter
