# envoy-cli

A lightweight CLI for managing environment variable sets across multiple deployment targets.

## Installation

```bash
go install github.com/your-org/envoy-cli@latest
```

## Commands

### `init`
Initialise a new `.envoy.json` config file in the current directory.
```bash
envoy init
```

### `target`
Manage deployment targets.
```bash
envoy target list
envoy target add <name>
envoy target remove <name>
envoy target rename <old> <new>
```

### `env`
Manage environment variables per target.
```bash
envoy env set <target> <key> <value>
envoy env get <target> <key>
envoy env delete <target> <key>
envoy env list <target>
```

### `diff`
Compare environment variables between two targets.
```bash
envoy diff <target1> <target2>
```

### `copy`
Copy or merge environment variables between targets.
```bash
envoy copy <src> <dest>
envoy copy --merge <src> <dest>
```

### `snapshot`
Save and restore point-in-time snapshots of a target's variables.
```bash
envoy snapshot create <target> <name>
envoy snapshot list <target>
envoy snapshot restore <target> <name>
```

### `export`
Export environment variables for a target in various formats.
```bash
# Print to stdout (default: dotenv)
envoy export <target>

# Choose format: dotenv | shell | json
envoy export <target> --format shell
envoy export <target> --format json

# Write directly to a file
envoy export <target> --format dotenv --out .env
```

## Config file

All data is stored in `.envoy.json` in the working directory.
