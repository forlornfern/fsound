package program

import (
	vlc "github.com/adrg/libvlc-go/v3"
	tea "github.com/charmbracelet/bubbletea"
)

type Fsound struct {
	PlaylistPaths []string `json:"playlist-paths"`
	program       *tea.Program
	playlists     []*playlist
	player        *vlc.Player
	width, height int
	err           error
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
		}
	case tea.WindowSizeMsg:
		f.width = msg.Width
		f.height = msg.Height
	}
	return f, nil
}

func (f *Fsound) View() string {
	return ""
}
