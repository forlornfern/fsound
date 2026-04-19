package cmd

import (
	"slices"

	"fsound/store"

	"charm.land/log/v2"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <directory>...",
	Short: "Add playlists",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if verbose, _ := rootCmd.PersistentFlags().GetBool("verbose"); verbose {
			logger.SetLevel(log.InfoLevel)
		}
		pm, err := store.LoadProgramModel()
		if err != nil {
			return err
		}
		for _, dir := range args {
			if isDir, err := afero.IsDir(afero.NewOsFs(), dir); err != nil || !isDir {
				logger.Errorf("%q is not a directory", dir)
				continue
			}
			if contains := slices.Contains(pm.PlaylistPaths, dir); contains {
				logger.Warnf("The %q directory has already been added", dir)
				continue
			}
			pm.PlaylistPaths = append(pm.PlaylistPaths, dir)
			logger.Infof("%q saved", dir)
		}

		return store.SaveProgramModel(pm)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
