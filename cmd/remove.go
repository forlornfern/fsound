package cmd

import (
	"fsound/store"
	"slices"

	"charm.land/log/v2"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove <directory>...",
	Aliases: []string{"rm"},
	Short:   "Remove playlists",
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
		pm, err := store.LoadProgramModel()
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		var completions []cobra.Completion
		for _, path := range pm.PlaylistPaths {
			if !slices.Contains(args, path) {
				completions = append(completions, path)
			}
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	},
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if verbose, _ := rootCmd.PersistentFlags().GetBool("verbose"); verbose {
			logger.SetLevel(log.InfoLevel)
		}
		pm, err := store.LoadProgramModel()
		if err != nil {
			return err
		}

		for _, dir := range args {
			i := slices.Index(pm.PlaylistPaths, dir)
			if i == -1 {
				logger.Errorf("%q is not directory", dir)
				continue
			}
			logger.Infof("%q removed", pm.PlaylistPaths[i])
			pm.PlaylistPaths = append(pm.PlaylistPaths[:i], pm.PlaylistPaths[i+1:]...)
		}

		return store.SaveProgramModel(pm)
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
