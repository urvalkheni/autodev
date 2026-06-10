# AutoDevs × DevMentor Scaffold

Copy this folder into any target repository to enable the voice-to-code supply chain.

## Quick setup

```bash
# From target repo root
cp -r path/to/templates/autodevs-scaffold/.autodevs ./
cp path/to/templates/autodevs-scaffold/AGENTS.md ./AGENTS.md
```

Edit `.autodevs/prompts-base.md` for your stack (Flutter, React, FastAPI, etc.).

## Codex

Open the repo in Codex. It reads `AGENTS.md` at root, which boots the Ralph Loop from `.autodevs/`.

```bash
cd /your/project
codex
```

Or run a specific todo:

```bash
codex "Execute Ralph Loop for .autodevs/todo/todo_001.md"
```

## DevMentor flow

1. Speak idea in DevMentor mobile
2. DevMentor writes `todo_YYYY_MM_DD_<slug>.md`
3. Commits to `/.autodevs/todo/` on selected GitHub repo
4. AutoDevs / Codex detects and runs Ralph Loop
5. PR opened → DevMentor notification

## Folder map

| Path | Purpose |
|------|---------|
| `todo/` | Incoming tasks from DevMentor |
| `completed/` | Finished tasks + reports |
| `failed/` | Failed after max loops |
| `loops/` | Per-iteration state |
| `memory/` | Learned patterns |
| `agents/` | Planner → Optimizer prompts |
