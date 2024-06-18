package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tsukinoko-kun/zzh/internal/config"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "create a new zzh config",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := config.New()
		return err
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
