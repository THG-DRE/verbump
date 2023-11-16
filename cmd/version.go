package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var Version = "undefined"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show the verbump version",
	Run:   versionCmdRunE,
}

func versionCmdRunE(cmd *cobra.Command, args []string) {
	fmt.Println(Version)
}
