#!/bin/bash

# Create the output directory if it doesn't exist
mkdir -p builds

# Build for all architectures
#Macos 
GOOS=darwin GOARCH=arm64 go build -o downloads/livecode_macos
#Linux
GOOS=linux GOARCH=amd64 go build -o downloads/livecode_linux
#Windows
GOOS=windows GOARCH=amd64 go build -o downloads/livecode_windows.exe