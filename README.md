# KODKAFA
> _A Persistent CLI for Your Scripts, With Memory_

```text
██╗  ███╗ ██████╗ ███████╗ ██╗  ███╗ ██████╗ ███████╗ ██████╗
██║ ███╔╝██╔═══██╗██╔═══██╗██║ ███╔╝██╔═══██╗██╔════╝██╔═══██╗
█████╔═╝ ██║   ██║██║   ██║█████╔═╝ ████████║█████╗  ████████║
██╔═███╗ ██║   ██║██║   ██║██╔═███╗ ██╔═══██║██╔══╝  ██╔═══██║
██║  ███╗╚██████╔╝███████╔╝██║  ███╗██║   ██║██║     ██║   ██║
╚═╝  ╚══╝ ╚═════╝ ╚══════╝ ╚═╝  ╚══╝╚═╝   ╚═╝╚═╝     ╚═╝   ╚═╝
```
[![Install with Homebrew](https://img.shields.io/badge/Homebrew-Install-2a7d2e?logo=homebrew&logoColor=white)](https://github.com/kodkafa/kod#macos-homebrew)
<a href="https://github.com/charmbracelet/bubbletea/releases"><img src="https://img.shields.io/github/release/kodkafa/kod" alt="Latest Release"></a>


KODKAFA is a persistent meta-runner for your local scripts and micro-tools. It manages them as **plugins**, provides a stunning TUI dashboard, and tracks **execution history per plugin** so you can resume workflows without re-typing complex flags.

The primary goals are:

*   **One entry point** for all dev micro-tools
*   **Per-plugin contextual history** (not global shell history)
*   **Predictable dependency management** (isolated environments via Core Architecture)
*   **A SOLID, testable core** with a clean separation between UI and Business Logic

---

## Quick Start

1.  **Initialize**: Run `kod` or `kod init` for the first time; it will create `~/.kodkafa/`.
2.  **Add a Sample**:
    *   Press `m` to open the Command Menu.
    *   Select `Add Plugin`.
    *   Enter the path to a sample: `./samples/hello-py` (or a URL).
3.  **Run**:
    *   Select the newly added `hello-py` from the Dashboard.
    *   Press `Enter` to run.
    *   Observe the "smart prompt" (it remembers your last args!).

## Usage

KODKAFA runs in two modes:

* **TUI mode (`kod`)**: interactive dashboard + smart prompt execution
* **CLI mode (`kod run ...`, `kod add ...`)**: direct commands without entering the TUI

### Commands

```text
kod                      # Open TUI Dashboard
kod init                 # Initialize ~/.kodkafa structure
kod list                 # List installed plugins
kod info <name>          # View plugin metadata & stats
kod add <path|url>       # Install a plugin
kod run <name>           # Execute a plugin directly
kod load <name>          # reload/install plugin dependencies
kod del <name>           # Remove a plugin
```

### Aliases

* `a` → `add`
* `l` → `load`
* `r` → `run`
* `i` → `info`
* `d` → `del`

---

## Configuration

Global configuration is stored in `~/.kodkafa/config.json`.

```json
{
    "splash": true,
    "items_per_page": 5,
    "sort_by": "recent",
    "show_last_runs": true,
    "supported_runtimes": {
        "python": "python3",
        "node": "node",
        "r": "Rscript"
    }
}
```

*   **splash**: Enable/Disable the startup ASCII art animation.
*   **items_per_page**: Number of plugins to show per page in the dashboard.
*   **supported_runtimes**: Customize the binary paths for different languages.

---

## Plugin Contract

Each plugin folder must include `plugin.yml`. KODKAFA supports Python, Node.js, R, and Shell scripts.

```yaml
name: my_tool
language: python  # Options: python, node, r, shell
description: A useful description for the dashboard
entry: run.py     # Main entry point file
usage: "args..."  # Hint for the prompt
```

> [PLUGIN.md](docs/PLUGIN.md)

### Supported Languages
*   **Python**: Runs in an isolated venv managed by KODKAFA.
*   **Node.js**: Runs with its own `node_modules` (handled via `pnpm`/`npm`).
*   **R**: Executed via `Rscript`.
*   **Shell**: Standard executable scripts.

---

## Architecture

KODKAFA follows a **SOLID** architecture, ensuring the UI is strictly a presentation layer.

*   **UI Layer (Bubble Tea)**: Handles rendering and user input.
*   **Application Layer**: Contains Use Cases (`AddPlugin`, `RunPlugin`, etc.) that orchestrate logic.
*   **Domain Layer**: Defines entities (`Plugin`, `RunRecord`) and interfaces (`Repository`, `Runner`).
*   **Infrastructure Layer**: Implements storage (JSON), process execution (`exec.Process`), and runtime management.

### Persistence Layout (`~/.kodkafa/`)
*   `plugins/` — Source code for installation plugins.
*   `state/` — Per-plugin execution history (`<plugin>.json`).
*   `core/` — Centralized runtime environments (e.g., Python venvs, Node modules).
*   `config.json` — User preferences.

---

## Installation (From Source)

1.  **Prerequisites**: Go 1.21+, Git.
2.  **Clone & Build**:
    ```bash
    git clone https://github.com/kodkafa/kodkafa-cli.git
    cd kodkafa-cli
    go mod download
    go build -o kod cmd/kod/main.go
    ```
3.  **Run**:
    ```bash
    ./kod
    ```

## Development

*   **Hot Reload**: Run `air` to start the dev server with hot reloading (configured via `.air.toml`).
*   **Testing**: Run `go test ./...` to execute core tests.

---

**License**: [MIT](LICENSE)
