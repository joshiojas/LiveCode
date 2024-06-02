#!/bin/bash

# Installation Script for LiveCode

# Functions

LIVE_CODE_BASE_URL="https://raw.githubusercontent.com/joshiojas/LiveCode/main/builds/"
BINARY_NAME="livecode"
INSTALL_DIR="${HOME}/bin"

# Functions

download_binary() {
    os_name="$(uname -s)"
    case "${os_name}" in
        Linux*)     os_ext="linux" ;;
        Darwin*)    os_ext="macos" ;;
        MINGW64*)   os_ext="windows"; BINARY_NAME="${BINARY_NAME}.exe" ;;
        *)          echo "Unsupported operating system: ${os_name}" && exit 1 ;;
    esac
    machine_type="$(uname -m)"

    BINARY_URL="${LIVE_CODE_BASE_URL}/${BINARY_NAME}_${os_ext}"

    echo "Downloading LiveCode for ${os_name} (${machine_type})..."
    curl -LO "${BINARY_URL}"
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
