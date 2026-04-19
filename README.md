# envoy-cli

A lightweight CLI for managing environment variable sets across multiple deployment targets.

## Installation

```bash
go install github.com/yourusername/envoy-cli@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/envoy-cli/releases).

## Usage

```bash
# Initialize a new envoy config in the current directory
envoy init

# Add an environment variable to a target
envoy set production DATABASE_URL=postgres://user:pass@host/db

# List all variables for a target
envoy list production

# Apply variables to a deployment target
envoy apply production

# Copy variables from one target to another
envoy copy staging production
```

### Example Config (`envoy.yaml`)

```yaml
targets:
  production:
    DATABASE_URL: postgres://user:pass@prod-host/db
    API_KEY: your-api-key
  staging:
    DATABASE_URL: postgres://user:pass@staging-host/db
    API_KEY: staging-api-key
```

## Commands

| Command | Description |
|---|---|
| `init` | Initialize a new envoy config |
| `set <target> <KEY=VALUE>` | Set a variable for a target |
| `list <target>` | List all variables for a target |
| `apply <target>` | Apply variables to a target |
| `copy <src> <dst>` | Copy variables between targets |

## Contributing

Pull requests are welcome. For major changes, please open an issue first.

## License

[MIT](LICENSE)