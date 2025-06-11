# Pirate Weather CLI

A minimal CLI for the [Pirate Weather API](https://pirateweather.net/)

## Build and Installation

This project uses [`just`](https://github.com/casey/just) as a command runner and requires [Go](https://golang.org/) for building. Below are the available commands:

### Available Commands

- `just` or `just install` - Build and install the binary to `PREFIX/bin/` (default: `/usr/local/bin`)
- `just build` - Build the binary in the current directory
- `just uninstall` - Remove the installed binary
- `just clean` - Remove the built binary from the current directory

## Usage

```sh
pirate
pirate -lat 40.7128 -lon -74.0060
pirate -units us
```

## Configuration

```
export PIRATE_WEATHER_API_KEY=your_api_key_here
export PIRATE_WEATHER_LAT=40.7128      # Default: New York City latitude
export PIRATE_WEATHER_LON=-74.0060     # Default: New York City longitude
export PIRATE_WEATHER_UNITS=us         # Options: us, si, ca, uk (Default: us)
```

![pirate](https://github.com/user-attachments/assets/88a88c1b-8b13-4371-b7b7-5c331c8496d8)
