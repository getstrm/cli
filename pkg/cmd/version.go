package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"pace/pace/pkg/common"
)

var VersionCmd = &cobra.Command{
	Use:               "version",
	DisableAutoGenTag: true,
	Short:             "Get the cli version",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(common.Version)
		return nil
	},
}
