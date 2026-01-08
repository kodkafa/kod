## KODKAFA — A Persistent CLI for Your Scripts, With Memory 

This document explains **what the program does** and describes its **end-to-end flows** in a way that can be directly used as a design spec for a SOLID-friendly Bubble Tea rewrite.

KODKAFA is a persistent meta-runner for local scripts and small CLI tools. It manages them as **plugins**, provides a single entry point (`kod`), and keeps **separate execution history per plugin** so you can resume workflows without re-remembering commands or flags.

KODKAFA supports two modes:

1. **TUI Mode (`kod`)**: plugin browsing, search/filter, plugin inspection, a “smart prompt” for reruns, and per-plugin history walking.
2. **CLI Mode (`kod run ...`, `kod add ...`, etc.)**: direct actions without opening the TUI, while still updating history and usage tracking.

---

## 1) Application Lifecycle

### 1.1 Startup

When the user runs `kod`, the program:

1. Verifies the presence of the `~/.kodkafa/` directory structure (or guides the user to initialize it).
2. Loads global configuration from `~/.kodkafa/config.json` (trusted domains, dependency strategy, TUI preferences, etc.).
3. Loads the plugin inventory (from `~/.kodkafa/plugins/` and/or an internal registry index).
4. Loads usage analytics from `~/.kodkafa/usage.json` (recently used, most used).
5. Starts the TUI and enters the **Dashboard (Plugin List)** screen.

### 1.2 Shutdown

On exit:

* Flushes any pending store writes (state/usage).
* Closes log streams and file handles.
* Ensures running processes (if any) are terminated or detached according to policy.

---

## 2) TUI Screens and States (State Machine)

The TUI is best modeled as a deterministic state machine. Each state defines:

* What is displayed
* Which inputs are handled
* Which transitions occur
* Which persistence writes (if any) happen

### UI Architecture (The Dual-List Model)

The dashboard is split into two distinct functional areas to balance speed and discoverability:

1.  **Top List (Favorites/Recents)**: A non-paginated, fixed section containing up to `Config.FavLimit` plugins. It follows the `Config.LastRunOrder` policy (either "most used" or "recently used").
2.  **Main Inventory (Bottom)**: A paginated, alphabetical list of all registered plugins, excluding those currently visible in the Top List.
3.  **Command Overlay**: Triggered by the `m` key, this modal-like menu allows for administrative tasks (Add, Delete, System Info) without leaving the TUI context.

**Branding & Feedback:**
*   **Header**: Features the animated ASCII logo, the tagline *"KODKAFA — A Persistent CLI for Your Scripts, With Memory"*, and the official URL.
*   **Footer**: A dynamic shortcut bar showing context-sensitive keybindings (e.g., `s:search`, `i:info`, `m:menu`).

---

### State A — Dashboard: Plugin List

**Inputs → Behavior:**

* `Up/Down`: change selected plugin.
* `s`: enter Search/Filter (State B).
* `i`: open Plugin Info (State C).
* `m`: open System Menu (State D).
* `Enter`: open Run Prompt for selected plugin (State E).
* `Esc`: exit (either immediate or via confirmation, depending on preferences).

**Persistence:** none.

---

### State B — Search / Filter

**Goal:** Filter the plugin list instantly.

**Inputs → Behavior:**

* Text input: updates filter string and narrows the list.
* `Enter`: apply filter and return to Dashboard (State A).
* `Esc`: cancel filter and return to Dashboard (State A).

**Persistence:** none.

---

### State C — Plugin Info

**Goal:** Show metadata and a preview of the plugin’s recent runs.

**UI:**

* Description
* Added date
* Last executed time
* Preview of last N parameter sets (for example, last 5)
* Optionally: runtime type, entrypoint, dependency status

**Inputs → Behavior:**

* `Esc` / `Left`: return to Dashboard (State A).
* Optional: `Enter` to go to Run Prompt (State E).

**Data sources:**

* `~/.kodkafa/state/<plugin>.json`
* `~/.kodkafa/usage.json`

---

### State D — System Menu

**Goal:** Provide access to system-level actions via a submenu (add/delete/load/init).

**Inputs → Behavior:**

* Select menu item + `Enter`: transition to the corresponding flow:

  * Add Flow (State D1)
  * Delete Flow (State D2)
  * Load Flow (State D3)
  * Init Flow (State D4)
* `Esc`: return to Dashboard (State A)

**Note:** These TUI flows should call the same application-layer use cases as CLI commands.

---

### State E — Run Prompt (“Smart Prompt”)

**Goal:** Capture parameters before running and enable per-plugin history navigation.

**UI:**

* Input line pre-filled with the **most recent parameters** for that plugin (if available).
* Hint text: “Up/Down for history”.

**Inputs → Behavior:**

* `Up`: move backward through the plugin’s parameter history; replace the input with the selected historical entry.
* `Down`: move forward toward the newest entry; optionally end in a “fresh prompt”.
* Text editing: behave like a typical single-line shell input (insert, delete, left/right, etc.).
* `Enter`: begin execution (State F).
* `Esc`: return to Dashboard (State A).

**Key behavior:**
The prompt represents `pluginName + argsString`. The program’s runner composes the actual command line based on runtime detection (Node/Python/R/Bash/etc.).

---

### State F — Running (Process View)

**Goal:** Execute the plugin and stream output live.

**On-enter actions:**

1. Parse the input into an argument string (or store raw input consistently).
2. Persist the new run intent (policy choice: record attempts always, or only successful runs):

   * Append a run entry to `~/.kodkafa/state/<plugin>.json` (timestamp, args; exit code can be updated after completion).
   * Update `~/.kodkafa/usage.json` (recent list and run counters).
3. Resolve runtime and build the final execution command.
4. Start the process and stream stdout/stderr to:

   * The TUI output view
   * A log file under `~/.kodkafa/logs/`

**Inputs → Behavior:**

* Output streams continuously while the process runs.
* On completion, transition to Post-Run (State G).

---

### State G — Post-Run (Return / Loop)

**Goal:** Let the user read output and decide to rerun or return.

**UI:**

* “Process finished.”
* Options: `[Enter] Run again` and `[Esc] Back to Menu`

**Inputs → Behavior:**

* `Enter`: return to Run Prompt (State E) (typically pre-filled with the most recent parameters).
* `Esc`: return to Dashboard (State A).

---

## 3) CLI Command Flows (Non-TUI)

### Command List

```text
kod
kod add <path|repo>
kod del <name>
kod load <name>
kod run <name>
kod info <name>
kod log <name>
kod init
```

---

### 3.1 `kod add <path|repo>`

**Goal:** Register a plugin from a local path or a trusted remote URL.

**Steps:**

1. Determine whether the input is a local path or a URL.
2. If URL:

   * Enforce trusted domain policy from `config.json`.
   * If not trusted, apply the program’s safety policy (reject by default, or require an explicit override).
3. Fetch/copy the plugin source into `~/.kodkafa/plugins/<name>/`.
4. Read and validate `plugin.yml` (name, description, runtime, entrypoint).
5. Update the plugin registry/index.
6. Create initial plugin state file `~/.kodkafa/state/<plugin>.json` with “Added Date” and empty history.

**Output:**

* “Plugin added: <name>”
* Optional: instruction to run `kod load <name>`.

---

### 3.2 `kod load <name>`

**Goal:** Install dependencies for a plugin.

**Steps:**

1. Locate plugin and read runtime metadata.
2. Choose dependency strategy:

   * Centralized runtime environments under `~/.kodkafa/core` (via Core Architecture), or
   * Isolated per-plugin environment if mandated by policy.
3. Detect dependency manifest (for example, `requirements.txt`, `package.json`).
4. Run the appropriate installer (e.g. `pip` for Python, `npm` for Node, `install.packages` for R).
5. Write detailed logs to `~/.kodkafa/logs/`.

**Output:**

* “Dependencies installed.” or a clear error report.

---

### 3.3 `kod del <name>`

**Goal:** Remove a plugin and optionally clean dependencies safely.

**Steps:**

1. Validate plugin existence.
2. Ask whether to remove dependencies (TUI confirmation, CLI flag policy).
3. Remove plugin source from `~/.kodkafa/plugins/<name>/`.
4. Remove plugin state and logs.
5. If dependency cleanup is selected:

   * Check whether other plugins share the same dependency resources.
   * Only reclaim disk space when safe.

**Output:**

* “Plugin removed: <name>” (+ dependency cleanup status)

---

### 3.4 `kod run <name> [args...]`

**Goal:** Run a plugin immediately without entering the TUI.

**Steps:**

1. Locate plugin and resolve runtime.
2. Persist run to per-plugin state and usage stats.
3. Execute the process and stream output to the terminal and logs.
4. Return the plugin’s exit code as the CLI exit status.

---

## 4) Persistence Contract (Data Writes)

KODKAFA’s “memory” is built on three persistent stores:

### 4.1 Per-Plugin State — `~/.kodkafa/state/<plugin>.json`

Stores:

* Added date
* Last executed time
* Run count
* Input history (last N runs, with timestamps and args)
* Optional: exit code, duration, status

**Write triggers:**

* `run` (TUI or CLI): append history entry, update last executed time, increment run count
* `add`: create initial state
* `del`: delete state file

---

### 4.2 Global Usage — `~/.kodkafa/usage.json`

Stores:

* Recently used plugins (last N, ordered by timestamp)
* Most used plugins (top N by run count)

**Write triggers:**

* Every `run` updates both recent ordering and usage counters.

---

### 4.3 Logs — `~/.kodkafa/logs/`

Stores:

* Plugin run logs (stdout/stderr)
* System operation logs (add/load/del/init)

**Write triggers:**

* `run`, `load`, and optionally other operations based on observability policy.

---

## 5) SOLID-Friendly Boundaries for a Bubble Tea Rewrite

The Bubble Tea layer should remain a **presentation layer**, while core behavior lives in application use cases. This reduces complexity and makes flows testable.

**Recommended responsibility split:**

### TUI Layer (Bubble Tea Model/Update/View)

* Handles input events
* Maintains UI state machine
* Calls application use cases
* Renders results

### Application Layer (Use Cases)

* `AddPlugin`
* `DeletePlugin`
* `LoadPluginDependencies`
* `RunPlugin`
* `GetPluginInfo`
* `ListPlugins`
* `GetPluginLog`

Each use case should have one clear responsibility and be unit-testable.

### Domain Layer

* Entities: `Plugin`, `PluginState`, `RunRecord`, `UsageStats`
* Interfaces (ports): `PluginRepository`, `StateStore`, `UsageStore`, `Runner`, `DependencyInstaller`, `Logger`

### Infrastructure Layer

* JSON file stores (read/write)
* Process execution (runner)
* Dependency installers (pip/npm/R wrappers)
* Logging implementation

This structure directly supports:

* Single Responsibility per module
* Dependency Inversion (TUI depends on abstractions; infrastructure implements them)
* Clean test seams for runners and stores

---

## 6) Concrete Example Flow: The `ss` Plugin

1. User runs `kod` → Dashboard (State A)
2. User highlights `ss` → `Enter` → Run Prompt (State E)
3. Prompt pre-fills with the most recent args for `ss`
4. User presses `Up` to navigate older runs; `Down` to return toward the newest
5. User edits the args → `Enter`
6. Running (State F):

   * `state/ss.json` and `usage.json` are updated
   * Process runs; output streams live; logs are written
7. Process completes → Post-Run (State G)
8. `Enter` reruns (back to State E), or `Esc` returns to Dashboard (State A)

**Outcome:**
The new command parameters become the latest entry in the `ss` history stack, ready for instant reuse next time.

---
