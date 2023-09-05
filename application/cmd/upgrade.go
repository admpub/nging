package cmd

import (
	"github.com/admpub/nging/v5/application/handler/setup"
	"github.com/admpub/nging/v5/application/library/config"
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade Database Struct",
	RunE:  upgradeRunE,
}

func upgradeRunE(cmd *cobra.Command, args []string) error {
	conf, err := config.InitConfig()
	config.MustOK(err)
	conf.AsDefault()
	return setup.Upgrade()
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}
