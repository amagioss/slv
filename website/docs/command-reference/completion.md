---
sidebar_position: 5
---

# Auto Complete
Documentation for the `slv completion` command.

## Set up Command Auto Complete
Hitting `tab` would auto complete the command for you within the `slv` command.

#### General Usage:
```bash
slv completion [command]
```

#### Commands Available:
- [`bash`](#bash)
- [`fish`](#fish)
- [`powershell`](#powershell)
- [`zsh`](#zsh)

### Bash

On bash versions > 4,
```bash
source <(slv completion bash)
```

You are likely to be using an older bash version only if you are using MacOS's default bash (3.2.x)
In that case, you can run
```bash
slv completion bash > /tmp/slvcomp.sh && source /tmp/slvcomp.sh && rm /tmp/slvcomp.sh
```

You can add either of the lines to `~/.bashrc` to automatically set up auto complete everytime.

---

### Fish
In Fish, the auto completion script can be added to the auto completions path.
```fish
slv completion fish > ~/.config/fish/completions/slv.fish
```
The completions can be loaded using the following command. Alternatively, you can spawn a new fish shell.
```fish
. ~/.config/fish/config.fish
```
---

### Powershell
To enable auto completion temporarily for the current session, 
```powershell
slv completion powershell | Out-String | Invoke-Expression
```

\
Alternatively, auto completion can be enabled persistently.
```powershell
slv completion powershell > $PROFILE.d/slv-completion.ps1
```
Add the following line to the PowerShell Profile
```powershell
Add-Content -Path $PROFILE -Value ". `$PROFILE.d/slv-completion.ps1"
```
Reload 
```powershell
. $PROFILE
```
---

### Zsh
The following commands will load the Zsh completion system. Ensure that these commands are either run or present on top of `~/.zshrc` before the auto completion script sourcing.
```zsh
autoload -Uz compinit
compinit
```

Now you can either run the following command or add it to `~/.zshrc` for SLV Auto Completion.
```zsh
source <(slv completion bash)
```



