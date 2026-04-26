package program

import (
	"strings"

	"charm.land/lipgloss/v2"
	vlc "github.com/adrg/libvlc-go/v3"
	tea "github.com/charmbracelet/bubbletea"
)

type Fsound struct {
	PlaylistPaths         []string `json:"playlist-paths"`
	program               *tea.Program
	playlists             []*playlist
	selectedPlaylistIndex int
	player                *vlc.Player
	width, height         int
	err                   error
}

func Execute(model *Fsound) (*Fsound, error) {
	p := tea.NewProgram(model)
	model.program = p
	m, err := p.Run()
	if err != nil {
		return nil, err
	}
	model = m.(*Fsound)
	if model.err != nil {
		return nil, model.err
	}

	return model, nil
}

func (f *Fsound) Init() tea.Cmd {
	for _, path := range f.PlaylistPaths {
		playlist, err := newPlaylist(path)
		if err != nil {
			f.err = err
			return tea.Quit
		}
		f.playlists = append(f.playlists, playlist)
	}
	f.player, f.err = vlc.NewPlayer()
	if f.err != nil {
		return tea.Quit
	}
	return nil
}

func (f *Fsound) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return f, tea.Quit
		case "up":
			f.selectedPlaylistIndex = max(0, f.selectedPlaylistIndex-1)
		case "down":
			f.selectedPlaylistIndex = min(len(f.playlists)-1, f.selectedPlaylistIndex+1)
		}
	case tea.WindowSizeMsg:
		f.width = msg.Width
		f.height = msg.Height
	}
	return f, nil
}

func (f *Fsound) View() string {
	themeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.White).
		Background(lipgloss.Black).
		BorderForeground(lipgloss.White).
		BorderBackground(lipgloss.Black)
	borderStyle := themeStyle.BorderStyle(lipgloss.RoundedBorder())

	var sb strings.Builder
	for i, pl := range f.playlists {
		if i == f.selectedPlaylistIndex {
			sb.WriteRune('>')
		}
		sb.WriteString(pl.name + "\n")
	}

	playlistPanel := borderStyle.Width(f.width / 3).Height(f.height).Render(sb.String())

	return playlistPanel
}
