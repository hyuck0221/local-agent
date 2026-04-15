# local-agent

Run a local LLM and plug it into [opencode](https://opencode.ai) with one command.

`local-agent` is a thin wrapper around [Ollama](https://ollama.com) that picks
a model for you, makes sure the server is running, and registers it as an
OpenAI-compatible provider in `opencode.json` — so you can type `opencode` and
start chatting with a fully local coder model.

## Install

**macOS / Linux**

```sh
curl -fsSL https://raw.githubusercontent.com/hyuck0221/local-agent/main/install.sh | sh
```

**Windows (PowerShell)**

```powershell
iwr -useb https://raw.githubusercontent.com/hyuck0221/local-agent/main/install.ps1 | iex
```

**Node users (no install)**

```sh
npx local-agent start
```

## Use

```sh
local-agent start                 # interactive model picker
local-agent start qwen2.5-coder:7b  # skip the picker
local-agent status                # show server + opencode wiring state
local-agent stop                  # stop the local server
local-agent models                # list installed + recommended models
```

After `start`, run `opencode` and pick `local-agent/<your model>` from the
model list.

## What it actually does

1. Installs Ollama if missing (`brew` / `winget` / official script).
2. Starts `ollama serve` in the background if it is not already listening.
3. Pulls the chosen model if you do not already have it.
4. Merges a `provider.local-agent` block into your opencode config, preserving
   every other provider and key:

   ```json
   {
     "provider": {
       "local-agent": {
         "npm": "@ai-sdk/openai-compatible",
         "name": "Local Agent (Ollama)",
         "options": { "baseURL": "http://localhost:11434/v1", "apiKey": "ollama" },
         "models": { "qwen2.5-coder:7b": { "name": "qwen2.5-coder:7b" } }
       }
     }
   }
   ```

Config path: `~/.config/opencode/opencode.json` (Unix) or
`%APPDATA%\opencode\opencode.json` (Windows).

## Build from source

```sh
go build -o local-agent .
./local-agent status
```

## License

MIT
