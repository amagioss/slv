---
sidebar_position: 4
---

# Interactive Mode (TUI)

SLV provides a Terminal User Interface (TUI) that offers an interactive, menu-driven way to manage your vaults, profiles, and environments without needing to remember CLI commands.

## Launching the TUI

#### General Usage:
```bash
slv tui
```

#### Alternative Commands:
The TUI can also be launched using any of these aliases:
- `slv ui`
- `slv interactive`
- `slv menu`
- `slv dashboard`

#### Example:
```bash
$ slv tui
```

The TUI will launch and display the main menu with options to manage vaults, profiles, and environments.

---

## Main Menu

When you launch the TUI, you'll see the main menu with the following options:

- **Vaults** - Manage and organize your vaults
- **Profiles** - View profile settings and environments
- **Environments** - Manage environments
- **Help** - View documentation and help

---

## Features

### Vault Management
The TUI provides a comprehensive interface for managing vaults:

- **Browse Vaults** - Navigate through your file system to find vault files
- **View Vault Details** - See vault information, secrets, and access control
- **Edit Vaults** - Modify vault configuration, add/remove secrets, and manage access
- **Create New Vaults** - Set up new vaults with an interactive wizard
- **Unlock Vaults** - Access encrypted vaults by entering your password

### Profile Management (Upcoming Feature)
- View active profile
- See profile settings
- Manage environments within profiles

### Environment Management
- View available environments
- See environment details
- Create new environments

---

## Use Cases

The TUI is particularly useful for:

- **New Users** - Provides a guided interface for learning SLV without memorizing commands
- **Interactive Workflows** - When you need to browse and explore vaults visually
- **Quick Operations** - Fast access to common tasks through keyboard shortcuts
- **Visual Feedback** - See vault contents, access lists, and configurations in a structured format

---

## Tips

- Use keyboard shortcuts for faster navigation
- The status bar at the bottom shows helpful hints for the current screen
- Press `Esc` multiple times to navigate back through the menu hierarchy

---

## See Also

- [Quick Start Guide](/docs/quick-start) - Get started with SLV
- [Vault Component](/docs/components/vault) - Learn more about vaults
- [Profile Component](/docs/components/profile) - Learn about profiles
- [Environment Component](/docs/components/environment) - Learn about environments
- [Command Reference](/docs/command-reference/vault/get) - Explore CLI commands as an alternative

