package store

import (
	"encoding/json"

	"fsound/program"

	"github.com/adrg/xdg"
	"github.com/spf13/afero"
)

func LoadProgramModel() (*program.Fsound, error) {
	fs := afero.NewOsFs()
	path, err := xdg.DataFile("fsound/data.json")
	if err != nil {
		return nil, err
	}
	exists, err := afero.Exists(fs, path)
	if err != nil {
		return nil, err
	} else if !exists {
		err := afero.WriteFile(fs, path, []byte("{}"), 0644)
		if err != nil {
			return nil, err
		}
		return &program.Fsound{}, nil
	}

	file, err := afero.ReadFile(fs, path)
	if err != nil {
		return nil, err
	}
	var data program.Fsound
	if err := json.Unmarshal(file, &data); err != nil {
		return nil, err
	}
	valid := []string{}
	for _, p := range data.PlaylistPaths {
		if exists, err := afero.Exists(fs, p); err != nil || !exists {
			continue
		}
		valid = append(valid, p)
	}
	data.PlaylistPaths = valid

	return &data, nil
}

func SaveProgramModel(pm *program.Fsound) error {
	if path, err := xdg.DataFile("fsound/data.json"); err != nil {
		return err
	} else {
		fs := afero.NewOsFs()
		file, err := json.MarshalIndent(pm, "", "\t")
		if err != nil {
			return err
		}

		return afero.WriteFile(fs, path, file, 0644)
	}
}
