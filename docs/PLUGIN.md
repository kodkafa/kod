# Plugin System Documentation

## Step by Step Guide

1. **Create a folder** with the name of your plugin. ex: `my-plugin`
2. **Create a `plugin.yml`** file in the `my-plugin` directory.

   ```yaml
   name: my-plugin
   description: "My awesome plugin"
   version: "1.0.0"
   interpreter: python
   entry: main.py
   ```

3. **Develop or copy your script** (e.g., `main.py`) in the `my-plugin` directory.
4. **Add your code as a plugin**:
   ```bash
   kod add ./my-plugin
   ```
5. **(Optional) Share with the public**: Push your code to GitHub and it becomes available for everyone.
   ```bash
   kod add https://github.com/[USER_NAME]/my-plugin.git
   ```

KODKAFA CLI loads plugins from `~/.kodkafa/plugins`. Each plugin is a directory containing a `plugin.yml` manifest and its source code.

**DO NOT** add/create manually your plugin in the `~/.kodkafa/plugins` directory. Use `kod add` command to add your plugin. 

## Configuration (`plugin.yml`)

The `plugin.yml` file is the entry point for KODKAFA to understand how to run your plugin.

| Field | Type | Description |
| :--- | :--- | :--- |
| `name` | `string` | **Required.** Unique identifier for the plugin. Used in `kod run <name>`. |
| `description` | `string` | A brief description of what the plugin does. |
| `version` | `string` | Version of the plugin (e.g., "1.0.0"). |
| `interpreter` | `string` | **Required.** The runtime to use. Supported: `python`, `node`, `r`. |
| `entry` | `string` | **Required.** The main file to execute (e.g., `main.py`, `index.js`). |
| `usage` | `string` | Example command for the user to see in help menus. |
| `args` | `list` | (Optional) List of arguments for documentation purposes. |

### Argument Definition (`args`)

Each item in `args` can have:
- `name`: Argument flag/name.
- `type`: Data type (e.g., `string`).
- `required`: `true` or `false`.

## Supported Runtimes & Isolation

KODKAFA enforces runtime isolation to prevent conflicts.

### Python (`interpreter: python`)

- **Isolation:** Uses a dedicated `venv` inside the plugin directory.
- **Dependencies:** `pip install` packages into this virtual environment.
- **Execution:** Runs with the virtual environment's python executable.

### Node.js (`interpreter: node`)

- **Isolation:** Uses `node_modules` inside the plugin directory.
- **Dependencies:** `npm install` packages locally.
- **Execution:** Runs with `node`.

### R (`interpreter: r`)

- **Isolation:** Uses `renv` library paths.
- **Dependencies:** `renv::restore()` is used to manage packages.
- **Execution:** Runs with `Rscript`.

## Examples

### Python Example

`plugin.yml`:

```yaml
name: hello-py
description: "A simple python hello world"
version: "1.0.0"
interpreter: "python"
entry: hello.py
usage: 'kod run hello-py --name "Antigravity"'
args:
  - name: name
    type: string
    required: true
```

### Node.js Example

`plugin.yml`:

```yaml
name: hello-node
description: "A simple node hello world"
version: "1.0.0"
interpreter: "node"
entry: index.js
usage: 'kod run hello-node --name "NodeUser"'
args:
  - name: name
    type: string
    required: true
```

### R Example

`plugin.yml`:

```yaml
name: hello-r
description: "A simple R hello world"
version: "1.0.0"
interpreter: r
entry: hello.R
usage: 'kod run hello-r --name "RUser"'
args:
  - name: name
    type: string
    required: true
```

## Directory Structure

Plugins are stored in:

```text
~/.kodkafa/plugins/
  ├── my-plugin/
  │   ├── plugin.yml
  │   ├── main.py
  │   └── requirements.txt
  └── another-plugin/
      ├── plugin.yml
      └── index.js
```