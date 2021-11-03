package ui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
    Quit    key.Binding
    Up      key.Binding
    Down    key.Binding
    Start   key.Binding
    Stop    key.Binding
    Restart key.Binding
    Remove  key.Binding
    Logs    key.Binding
    Filter  key.Binding
    Refresh key.Binding
    Help    key.Binding
    Back    key.Binding
    Enter   key.Binding
}

var Keys = keyMap{
    Quit: key.NewBinding(
        key.WithKeys("q", "ctrl+c"),
        key.WithHelp("q", "quit"),
    ),
    Up: key.NewBinding(
        key.WithKeys("up", "k"),
        key.WithHelp("↑/k", "up"),
    ),
    Down: key.NewBinding(
        key.WithKeys("down", "j"),
        key.WithHelp("↓/j", "down"),
    ),
    Start: key.NewBinding(
        key.WithKeys("s"),
        key.WithHelp("s", "start"),
    ),
    Stop: key.NewBinding(
        key.WithKeys("t"), // 's' is taken, using 't' for stop
        key.WithHelp("t", "stop"),
    ),
    Restart: key.NewBinding(
        key.WithKeys("r"),
        key.WithHelp("r", "restart"),
    ),
    Remove: key.NewBinding(
        key.WithKeys("d"),
        key.WithHelp("d", "remove"),
    ),
    Logs: key.NewBinding(
        key.WithKeys("l"),
        key.WithHelp("l", "logs"),
    ),
    Filter: key.NewBinding(
        key.WithKeys("f"),
        key.WithHelp("f", "filter"),
    ),
    Refresh: key.NewBinding(
        key.WithKeys("f5"),
        key.WithHelp("F5", "refresh"),
    ),
    Help: key.NewBinding(
        key.WithKeys("?"),
        key.WithHelp("?", "help"),
    ),
    Back: key.NewBinding(
        key.WithKeys("esc"),
        key.WithHelp("esc", "back"),
    ),
    Enter: key.NewBinding(
        key.WithKeys("enter"),
        key.WithHelp("enter", "enter"),
    ),
}
