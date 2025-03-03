# CDX - Directory Alias Manager

CDX is a simple CLI tool that helps you manage directory aliases for quick navigation in Windows. Set aliases for your frequently used directories and open them in File Explorer with a single command.

## Features

- Set aliases for directories
- Remove aliases
- List all saved aliases
- Open File Explorer directly to aliased locations

## Installation

1. Download the latest release from GitHub
2. Place the `cdx.exe` in a directory that's in your PATH

## Usage

```bash
# Set an alias for a directory
cdx set desktop C:\Users\YourUsername\Desktop

# Remove an alias
cdx remove desktop

# List all saved aliases
cdx list

# Open File Explorer at aliased location
cdx desktop
```

## Storage

Aliases are stored in `%USERPROFILE%\.cdx_aliases.json`

## Building from Source

1. Clone the repository
2. Make sure you have Go installed
3. Run `go build`

## License
[MIT License](LICENSE)