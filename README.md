# docker-manager
# 1. Prerequisites

Install Go (1.21+)
Install Docker

# 2. Build the Application
Clone the project 

Initialize and build
cd docker-manager
go mod tidy
go build -o docker-manager .

# 3. Install (Optional)

sudo cp docker-manager /usr/local/bin/

# 4. Verify Docker Connection

Ensure Docker daemon is running
docker ps
# Example Usage

**Launch interactive TUI mode:**

./docker-manager interactive

**Launch with compact view:**

./docker-manager interactive --compact

**List containers in static mode:**

./docker-manager list

./docker-manager list --all

**Show real-time stats:**

./docker-manager stats

**View container logs:**

./docker-manager logs my-container

./docker-manager logs --tail 50 my-container

# Key Features

Interactive TUI: Full Bubbletea-based interface with keyboard controls

Container Management: Start, stop, restart, remove containers

Real-time Monitoring: Live CPU, memory, and network statistics

Logs Viewer: Scrollable logs display for selected containers

Filtering: Filter containers by name, status, or image

Compact Mode: Simplified view for smaller terminals

Static Commands: Non-interactive commands for scripting

Color-coded Status: Green for running, red for stopped, yellow for warnings

Resource Thresholds: Color-coded CPU and memory usage

Error Handling: Graceful handling of Docker connection issues

# Keyboard Shortcuts (Interactive Mode)
↑/↓: Navigate containers

s: Start container

t: Stop container

r: Restart container

d: Remove container

l: View logs

f: Filter containers

F5: Refresh

q: Quit

esc: Back
