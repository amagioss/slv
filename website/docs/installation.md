---
sidebar_position: 2
---
# Installation

## Installation

### MacOS (Homebrew)
You can install SLV in mac using [homebrew](https://brew.sh/), the package manager for MacOS.
```bash
brew install amagioss/slv/slv
```

### Linux
You can use the official installation script to install SLV on linux. (This will work on MacOS too.)
```bash
curl -fsSL https://slv.sh/scripts/install.sh | sh
```

### Windows
You can use the powershell script to install SLV on windows
```powershell
irm https://slv.sh/scripts/install.ps1 | iex
```

---

## Install a specific version
You can use the same install script to install specific versions of SLV 
### MacOS/Linux
```bash
curl -fsSL https://slv.sh/scripts/install.sh | sh -s v0.16.3
```
### Windows (Powershell)
```powershell
$v="v0.16.3"; irm https://slv.sh/scripts/install.ps1 | iex
```

---

## Verify Installation
Once installed, you can verify if SLV was installed correctly by entering 
```bash
slv --version
```
You should see the SLV version and build information in the console.

---
