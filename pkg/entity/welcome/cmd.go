package welcome

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"pace/pace/pkg/common"
	"path"
)

func DisableCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "welcome",
		Short:             "Disable welcoming message",
		DisableAutoGenTag: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			configPath, err := common.ConfigPath()
			if err != nil {
				return err
			}

			lastSeenCommandFilepath := path.Join(configPath, common.DefaultLastSeenFilename)
			_ = os.WriteFile(
				lastSeenCommandFilepath,
				[]byte("9999999999"),
				0644,
			)
			fmt.Println("Welcome message disabled")
			return nil
		},
		Args: cobra.ExactArgs(0),
	}
	return cmd
}
