package cmd

import (
	"os"

	"fsound/program"
	"fsound/store"

	"charm.land/log/v2"
	"github.com/spf13/cobra"
)

var (
	logger  = log.NewWithOptions(os.Stderr, log.Options{Level: log.WarnLevel, ReportTimestamp: false})
	rootCmd = &cobra.Command{
		Use:           "fsound",
		Short:         "Run fsound",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			pm, err := store.LoadProgramModel()
			if err != nil {
				return err
			}
			pm, err = program.Execute(pm)
			if err != nil {
				return err
			}
			store.SaveProgramModel(pm)

			return nil
		},
	}
)

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose log")
}

func Exec() {
	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(err)
	}
}
