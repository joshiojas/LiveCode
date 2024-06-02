#!/bin/bash

# Installation Script for LiveCode

# Constants
LIVE_CODE_URL="https://github.com/joshiojas/LiveCode/releases/download/macos/LiveCode"
BINARY_NAME="LiveCode"
INSTALL_DIR="${HOME}/.local/bin"

# Functions

download_binary() {
    echo "Downloading LiveCode..."
    curl -LO "${LIVE_CODE_URL}"
    chmod +x "${BINARY_NAME}"
}

install_binary() {
    echo "Installing LiveCode to ${INSTALL_DIR}..."
    mkdir -p "${INSTALL_DIR}"
    mv "${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
}

add_to_path() {
    if [[ ! "$PATH" =~ (^|:)${INSTALL_DIR}(:|$|) ]]; then
        echo "Adding ${INSTALL_DIR} to PATH..."
        echo "export PATH=\"${PATH}:${INSTALL_DIR}\"" >> "${HOME}/.zshrc"  # Add quotes
    fi
}

# Main Script

if ! command -v curl &> /dev/null; then
    echo "Error: curl is required to download LiveCode. Please install it."
    exit 1
fi

if [[ ! -w "${HOME}/.bashrc" && ! -w "${HOME}/.zshrc" ]]; then
    echo "Error: Cannot write to either .bashrc or .zshrc. Please adjust permissions or run as sudo."
    exit 1
fi

download_binary
install_binary
add_to_path

echo "Installation complete! LiveCode is now available in your PATH."
echo "Please restart your terminal session or source your shell profile to use it."
