---
sidebar_position: 6
---

# Help

Documentation for the `slv help` command.

## Getting help in SLV

### The `help` command
Get help for any command within SLV. 

####  Usage:
```bash
slv help <PATH_TO_COMMAND>
```
#### Example:
```bash
$ slv help env add
Adds an environment to the current profile

Usage:
  slv env add [flags]

Aliases:
  add, set, put, store, a

Flags:
  -e, --env-def strings   Environment definition
  -h, --help              help for add
      --root              Set the given environment as root
```

### The `--help` flag
Alternatively, all commands have a --help flag associated with it. You can simply type the command followed by `--help` to get the help text.

#### Example:
```bash
$ slv profile new --help
Creates a new profile

Usage:
  slv profile new [flags]

Flags:
      --git string          Git URI to clone the profile from
      --git-branch string   Git branch corresponding to the git URI
  -h, --help                help for new
  -n, --name string         Profile name
```
