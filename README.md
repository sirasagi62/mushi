# mushi - gitignore Template Generator

`mushi` is a command-line tool for generating `.gitignore` files from templates, inspired by the [github/gitignore](https://github.com/github/gitignore) repository. It provides both interactive and non-interactive modes to quickly create comprehensive ignore files for your projects.

## Features

- **Template-based generation**: Use templates from the official github/gitignore repository
- **Interactive mode**: Fuzzy-search and select templates with a TUI interface
- **Customizable common.gitignores**: Define your own default ignore patterns
- **Local caching**: Templates are cached locally for fast access
- **Force overwrite**: Option to overwrite existing `.gitignore` files

## Installation

```bash
go install github.com/sirasagi62/mushi@latest
```

## Usage

### Basic Usage

Generate a `.gitignore` file using a specific template:

```bash
mushi create Go
```

### Interactive Mode

Select a template interactively with fuzzy search:

```bash
mushi create -i
```

### Force Overwrite

Overwrite an existing `.gitignore` file:

```bash
mushi create Go -f
# or
mushi create Go --force
```

### Cache Management

Update the local template cache:

```bash
mushi cache update
```

Clean the local cache:

```bash
mushi cache clean
```

## Configuration

`mushi` uses the following directories and files:

- **Configuration**: `~/.config/mushi/`
- **Common ignore file**: `~/.config/mushi/common.gitignore` (automatically created with default patterns)
- **Cache directory**: `~/.cache/mushi/github-gitignore/` (local clone of github/gitignore)

The `common.gitignore` file contains default ignore patterns that are prepended to every generated `.gitignore` file. You can edit this file to customize your default ignores.

## How It Works

1. On first run, `mushi` clones the [github/gitignore](https://github.com/github/gitignore) repository to your local cache
2. When you create a `.gitignore`, it:
   - Updates the local cache (unless disabled)
   - Reads your custom `common.gitignore` file
   - Combines it with the selected template
   - Writes the result to `./.gitignore`

## License

The program itself, written in Golang, is distributed under the MIT License. See the [LICENSE](LICENSE) file for details.
However, `cmd/default.txt` and the default `~/.config/mushi/common.gitignore` generated from it are distributed under CC0.
See [./third_party_licenses](/third_party_licenses) for licenses of dependent third-party libraries.
