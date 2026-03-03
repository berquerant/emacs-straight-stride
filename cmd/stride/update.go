package main

import (
	"errors"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().Bool("all", false, "update all packages")
	updateCmd.Flags().BoolP("commit", "c", false, "update package and commit")
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

		if err := s.Update(cmd.Context(), pkg); err != nil {
			return err
		}
		if commit, _ := cmd.Flags().GetBool("commit"); commit {
			return s.Commit(cmd.Context())
		}
		return nil
	},
}
