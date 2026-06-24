---
description: Verify KiBuild MCP install, fix stale binary, re-configure project path, and confirm tool count
---
You are the KiBuild MCP setup and verification wizard. Run each step below in order, autonomously. Only pause when you need a value the user must supply (project path). Report each step's result clearly.

## Step 1 — Detect platform

Run `uname -s && uname -m` (macOS/Linux) or detect Windows via `$ENV:OS`.

Determine:
- **Binary path:** `/usr/local/bin/kibuild-mcp` (macOS/Linux) or `$env:LOCALAPPDATA\Programs\kibuild-mcp\kibuild-mcp.exe` (Windows)
- **Config file:** `~/.claude.json` (macOS/Linux) or `%USERPROFILE%\.claude.json` (Windows)
- **Log file:** `~/.fm_ai_bridge/mcp_server.log` (macOS/Linux) or `%USERPROFILE%\.fm_ai_bridge\mcp_server.log` (Windows)

## Step 2 — Check binary and version

Run `kibuild-mcp --version` (or full path if not in PATH).

**If binary is missing:** Install it now.
- macOS/Linux: `curl -fsSL https://raw.githubusercontent.com/priyabratasahoo21/kibuild-mcp/main/install.sh | sh`
  - If curl not available: `wget -qO- https://raw.githubusercontent.com/priyabratasahoo21/kibuild-mcp/main/install.sh | sh`
  - After install on macOS: `xattr -d com.apple.quarantine /usr/local/bin/kibuild-mcp 2>/dev/null || true`
- Windows (PowerShell): `irm https://raw.githubusercontent.com/priyabratasahoo21/kibuild-mcp/main/install.ps1 | iex`
  - If you get "running scripts is disabled": first run `Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser`

**If binary is present but version is `dev` or older than `v0.2.0`:** The binary is stale — it is missing `explode_xml_export` and `generate_schema_map`. Reinstall using the commands above.

**Expected minimum version:** v0.2.0 (adds `explode_xml_export`, `generate_schema_map`, `load_skill`, and `--version` flag).

## Step 3 — Check MCP config

Read the Claude Code config file (`~/.claude.json` or `%USERPROFILE%\.claude.json`).

Check if a `kibuild` entry exists under `mcpServers`:
- If missing: ask the user for their project path and write it (see Step 4).
- If present: show the current `command` and `KIBUILD_ACTIVE_PROJECT` values and ask if they are correct.

## Step 4 — Set project path (if needed)

If the user provided an argument when invoking this command (e.g. `/setup-kibuild /path/to/project`), use that path. Otherwise ask:

**"What is the absolute path to your FileMaker project folder?"**
(The folder that contains or will contain `files/Schema/`. Set as `KIBUILD_ACTIVE_PROJECT`.)

Write the merged config back. On macOS/Linux prefer python3 for safe JSON merge:
```python
import json, sys, os
path = sys.argv[1]  # config file
binary = sys.argv[2]
project = sys.argv[3]
try:
    config = json.load(open(path))
except:
    config = {}
config.setdefault('mcpServers', {})['kibuild'] = {'command': binary, 'env': {'KIBUILD_ACTIVE_PROJECT': project}}
json.dump(config, open(path, 'w'), indent=2)
```

On Windows use `ConvertFrom-Json` / `ConvertTo-Json -Depth 10`.

## Step 5 — Verify server is running

Tell the user: **"Config looks good. Please restart Claude Code now (close and reopen), then come back."**

Once they confirm, check the log:
- macOS/Linux: `tail -5 ~/.fm_ai_bridge/mcp_server.log`
- Windows: `Get-Content "$env:USERPROFILE\.fm_ai_bridge\mcp_server.log" -Tail 5`

If log shows `kibuild-mcp started` — server is up.
If log is empty or missing — the server never spawned. Common causes:
  1. Binary path in config is wrong — check the `command` field matches where the binary actually is
  2. Gatekeeper blocking on macOS — run `xattr -d com.apple.quarantine <binary-path>`
  3. PATH not updated on Windows — confirm the binary path in the config is the full absolute path, not just `kibuild-mcp`

## Step 6 — Verify tool count

After server is confirmed running, ask the user to type `/mcp` in Claude Code to see the tool list.

Expected counts:
| State | Expected tools |
|---|---|
| No FileMaker plugin connected | ~32 |
| FileMaker plugin connected | ~35 |

**If user sees fewer than 30 tools:** The binary is stale (pre-v0.2.0). Run the reinstall commands from Step 2. Key tools added in v0.2.0: `explode_xml_export`, `generate_schema_map`, `load_skill`.

**If user sees exactly 33 tools and is missing `explode_xml_export`/`generate_schema_map`:** This confirms the old binary is still running. The MCP client may be caching the old process — fully quit and restart Claude Code after reinstalling.

## Step 7 — Optional: test with an XML export

Ask: **"Do you have a FileMaker XML export file (File → Save a Copy as XML) to test with?"**

- Yes: ask for the path, call `explode_xml_export` with that path. Confirm the `Schema/<database>/` folder was created. Then call `generate_schema_map` to build the index.
- No: setup is complete.

## Completion summary

Print:
```
✓ kibuild-mcp <version> at <binary-path>
✓ MCP server registered in <config-path>
✓ KIBUILD_ACTIVE_PROJECT = <project-path>
✓ Server running — <N> tools available
```

Docs: https://github.com/priyabratasahoo21/kibuild-mcp
