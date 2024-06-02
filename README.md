# LiveCode - Automatic Command Restarter with PTY

LiveCode is a command-line tool written in Go that allows you to run a command within a pseudo-terminal (PTY) and automatically restart it whenever file changes are detected in the current directory (or optionally, in specified subdirectories). This is particularly useful for development workflows where you want to see the effects of code changes immediately.

## Features

- **PTY Support:** Runs commands in a PTY to capture colored output and maintain interactive features.
- **Automatic Restart on File Changes:** Restarts the command whenever files are created, modified, renamed, or deleted in the watched directory (or subdirectories).
- **Graceful Shutdown:** Attempts to gracefully terminate the running command before restarting.
- **Configurable:** You can specify the command to run and its arguments through command-line flags.
- **Recursive Watching (Optional):** You can choose to watch subdirectories for changes as well.

## Installation

### Pre-built Binaries

```bash
curl -fsSL https://raw.githubusercontent.com/joshiojas/LiveCode/main/install.sh | bash
```

### Deletion of Pre-built Binaries

```bash
curl -fsSL https://raw.githubusercontent.com/joshiojas/LiveCode/main/uninstall.sh | bash
```

## Building From Source

1. **Prerequisites:**

   - Go (version 1.22 or later) must be installed. You can download it from the [official website](https://golang.org/dl/).

2. **Building:**
   - Clone the repository:
     ```bash
     git clone https://github.com/joshiojas/LiveCode.git
     ```
   - Navigate to the project directory and build the executable:
     ```bash
     cd LiveCode
     go get .
     go build -o LiveCode
     ```
   - (Optional) Move the `LiveCode` executable to a directory in your `PATH` for easy access.

## Usage

```bash
livecode -cmd "<your-command>" [arguments]
```
