package main

import (
	"errors"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().Bool("all", false, "update all packages")
	updateCmd.Flags().BoolP("commit", "c", false, "update packages and commit")
}

var updateCmd = &cobra.Command{
	Use:   "update [PACKAGE...]",
	Short: `Update packages`,
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		if !all && len(args) == 0 {
			return errors.New("--all or PACKAGE is required")
		}

		pkgs := []string{}
		if !all {
			pkgs = args
		}
		s, err := newStraight(cmd)
		if err != nil {
			return err
		}
		defer s.Close()

		if err := s.Update(cmd.Context(), pkgs...); err != nil {
			return err
		}
		if commit, _ := cmd.Flags().GetBool("commit"); commit {
			return s.Commit(cmd.Context())
		}
		return nil
	},
}
