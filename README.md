# docker-manager

A terminal UI for managing Docker containers — built with [Bubbletea](https://github.com/charmbracelet/bubbletea) and [Cobra](https://github.com/spf13/cobra).

## Features

- Live container table with auto-refresh (every 5s)
- Start / stop containers with a keypress
- Stream logs into the TUI viewport
- Filter containers by name
- CPU & memory stats

## Install

```bash
go install github.com/farhanahmed3223/docker-manager@latest
```

## Usage

```bash
# Interactive TUI
docker-manager interactive

# CLI commands
docker-manager list
docker-manager stats <container-name>
docker-manager logs <container-name> --tail 50 --follow
```

## Keys

| Key | Action |
|-----|--------|
| ↑/↓ or j/k | Navigate |
| s | Start selected container |
| x | Stop selected container |
| l | View logs |
| / | Filter |
| r | Refresh |
| q | Quit |
