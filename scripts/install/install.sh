#!/bin/sh

set -e

if ! command -v unzip >/dev/null && ! command -v 7z >/dev/null; then
	echo "Error: either unzip or 7z is required to install SLV" 1>&2
	exit 1
fi

if [ "$OS" = "Windows_NT" ]; then
	target="windows-amd64"
else
	case $(uname -sm) in
	"Darwin x86_64") target="darwin_amd64" ;;
	"Darwin arm64") target="darwin_arm64" ;;
	"Linux aarch64") target="linux_arm64" ;;
	*) target="linux_amd64" ;;
	esac
fi

if [ $# -eq 0 ]; then
	slv_uri="https://github.com/amagioss/slv/releases/latest/download/slv_${target}.zip"
else
	slv_uri="https://github.com/amagioss/slv/releases/download/${1}/slv_${target}.zip"
fi

slv_install="${SLV_INSTALL:-$HOME/.slv}"
bin_dir="$slv_install/bin"
exe="$bin_dir/slv"

if [ ! -d "$bin_dir" ]; then
	mkdir -p "$bin_dir"
fi

curl --fail --location --progress-bar --output "$exe.zip" "$slv_uri"
if command -v unzip >/dev/null; then
	unzip -d "$bin_dir" -o "$exe.zip"
else
	7z x -o"$bin_dir" -y "$exe.zip"
fi
chmod +x "$exe"
rm "$exe.zip"

echo "SLV was installed successfully to $exe"
if command -v slv >/dev/null; then
	echo "Run 'slv --help' to get started"
else
	case $SHELL in
	/bin/zsh) shell_profile=".zshrc" ;;
	*) shell_profile=".bashrc" ;;
	esac
	echo "Manually add the directory to your \$HOME/$shell_profile (or similar)"
	echo "  export SLV_INSTALL=\"$slv_install\""
	echo "  export PATH=\"\$SLV_INSTALL/bin:\$PATH\""
	echo "Run '$exe --help' to get started"
fi
echo
