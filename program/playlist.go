package program

import (
	"path/filepath"
	"strings"

	vlc "github.com/adrg/libvlc-go/v3"
	"github.com/spf13/afero"
)

type playlist struct {
	name      string
	location  string
	mediaList []*vlc.Media
}

var audioExts = map[string]bool{
	".mp3":  true,
	".flac": true,
	".ogg":  true,
	".opus": true,
	".wav":  true,
	".aac":  true,
	".m4a":  true,
	".wma":  true,
	".aiff": true,
	".ape":  true,
}

func newPlaylist(directory string) (*playlist, error) {
	fs := afero.NewOsFs()
	files, err := afero.ReadDir(fs, directory)
	if err != nil {
		return nil, err
	}
	var mediaList []*vlc.Media
	for _, p := range files {
		if !p.IsDir() {
			path := filepath.Join(directory, p.Name())
			if isAudio(path) {
				if media, err := vlc.NewMediaFromPath(path); err != nil {
					continue
				} else {
					mediaList = append(mediaList, media)
				}
			}
		}
	}
	return &playlist{
		name:      filepath.Base(directory),
		location:  directory,
		mediaList: mediaList,
	}, nil
}

func isAudio(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return audioExts[ext]
}
