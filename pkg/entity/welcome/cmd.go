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
		Run: func(cmd *cobra.Command, args []string) {
			lastSeenCommandFilepath := path.Join(common.ConfigPath(), common.DefaultLastSeenFilename)
			os.WriteFile(
				lastSeenCommandFilepath,
				[]byte(fmt.Sprintf("%d", 9999999999)),
				0644,
			)
		},
		Args: cobra.ExactArgs(0),
	}
	return cmd
}
