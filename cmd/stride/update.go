package main

import (
	"errors"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().Bool("all", false, "update all packages")
}

var updateCmd = &cobra.Command{
	Use:   "update [PACKAGE]",
	Short: `Update package`,
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if !all && len(args) == 0 {
			return errors.New("--all or PACKAGE is required")
		}

		pkg := ""
		if !all {
			pkg = args[0]
		}
		s, err := newStraight(cmd)
		if err != nil {
			return err
		}
		defer s.Close()

		return s.Update(cmd.Context(), pkg)
	},
}
