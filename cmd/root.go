package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "liab",
	Short: "LIAB is a lightweight CLI for local-first and private software",
	Long:  `Local In A Box (LIAB) helps developers scaffold, secure, and manage offline-first architectures instantly.`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
