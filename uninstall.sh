#!/bin/bash

# Uninstallation Script for LiveCode

# Constants
BINARY_NAME="LiveCode"
INSTALL_DIR="${HOME}/.local/bin"

# Functions

remove_binary() {
    echo "Removing LiveCode from ${INSTALL_DIR}..."
    rm -f "${INSTALL_DIR}/${BINARY_NAME}"
}

remove_from_path() {
    if [[ "$PATH" =~ (^|:)${INSTALL_DIR}(:|$|) ]]; then
        echo "Removing ${INSTALL_DIR} from PATH..."
        sed -i '' "/export PATH=\"${PATH}:${INSTALL_DIR}\"/d" "${HOME}/.zshrc"
    fi
}

# Main Script

if [[ ! -w "${HOME}/.bashrc" && ! -w "${HOME}/.zshrc" ]]; then
    echo "Error: Cannot write to either .bashrc or .zshrc. Please adjust permissions or run as sudo."
    exit 1
fi

remove_binary
remove_from_path

echo "Uninstallation complete! LiveCode has been removed from your PATH."
echo "Please restart your terminal session or source your shell profile to apply changes."