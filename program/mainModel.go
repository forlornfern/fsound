package program

import (
	"path/filepath"
	"slices"
	"strings"

	"charm.land/lipgloss/v2"
	vlc "github.com/adrg/libvlc-go/v3"
	tea "github.com/charmbracelet/bubbletea"
)

type Fsound struct {
	PlaylistPaths    []string `json:"playlist-paths"`
	program          *tea.Program
	player           *vlc.Player
	playlists        []*playlist
	selectedPlaylist *playlist
	focusedPanel     int
	width, height    int
	err              error
}

const (
	playlistPanel int = iota
	mediaListPanel
)

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
	f.focusedPanel = mediaListPanel
	for _, path := range f.PlaylistPaths {
		playlist, err := newPlaylist(path)
		if err != nil {
			f.err = err
			return tea.Quit
		}
		f.playlists = append(f.playlists, playlist)
	}
	f.selectedPlaylist = f.playlists[0]
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
			switch f.focusedPanel {
			case playlistPanel:
				f.selectedPlaylist = f.playlists[max(0, slices.Index(f.playlists, f.selectedPlaylist)-1)]
			case mediaListPanel:
				f.selectedPlaylist.selectedMediaIndex = max(0, f.selectedPlaylist.selectedMediaIndex-1)
			}
		case "down":
			switch f.focusedPanel {
			case playlistPanel:
				f.selectedPlaylist = f.playlists[min(len(f.playlists)-1, slices.Index(f.playlists, f.selectedPlaylist)+1)]
			case mediaListPanel:
				f.selectedPlaylist.selectedMediaIndex = min(len(f.selectedPlaylist.mediaList)-1, f.selectedPlaylist.selectedMediaIndex+1)
			}
		case "left":
			f.focusedPanel = playlistPanel
		case "right":
			f.focusedPanel = mediaListPanel
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
		BorderForeground(lipgloss.BrightBlack).
		BorderBackground(lipgloss.Black)
	borderStyle := themeStyle.BorderStyle(lipgloss.RoundedBorder())
	playlistStyle := borderStyle.Width(f.width / 3).Height(f.height)
	medialistStyle := borderStyle.Width(f.width - playlistStyle.GetWidth()).Height(f.height)

	switch f.focusedPanel {
	case playlistPanel:
		playlistStyle = playlistStyle.BorderForeground(lipgloss.White)
	case mediaListPanel:
		medialistStyle = medialistStyle.BorderForeground(lipgloss.White)
	}

	var sb strings.Builder
	for i, pl := range f.playlists {
		if i == slices.Index(f.playlists, f.selectedPlaylist) {
			sb.WriteRune('>')
		} else {
			sb.WriteRune(' ')
		}
		sb.WriteString(pl.name + "\n")
	}

	leftPanel := playlistStyle.Render(sb.String())
	sb.Reset()

	for i, media := range f.selectedPlaylist.mediaList {
		if i == f.selectedPlaylist.selectedMediaIndex {
			sb.WriteRune('>')
		} else {
			sb.WriteRune(' ')
		}
		path, err := media.Location()
		ext := filepath.Ext(path)
		if err != nil {
			f.err = err
			f.program.Quit()
		}
		sb.WriteString(strings.TrimSuffix(filepath.Base(path), ext) + "\n")
	}
	rightPanel := medialistStyle.Render(sb.String())

	return lipgloss.JoinHorizontal(lipgloss.Center, leftPanel, rightPanel)
}
