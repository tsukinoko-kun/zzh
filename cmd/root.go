package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tsukinoko-kun/zzh/internal/config"
	"github.com/tsukinoko-kun/zzh/internal/ssh"
)

var rootCmd = &cobra.Command{
	Use:   "zzh",
	Short: "connect to a ssh server",
	RunE: func(cmd *cobra.Command, args []string) error {
		configs, err := config.GetConfigs()
		if err != nil {
			return err
		}
		config, err := config.Select(configs)
		if err != nil {
			return err
		}
		return ssh.Interactive(config)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
