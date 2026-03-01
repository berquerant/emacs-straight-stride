package main

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(rollbackCmd)
}

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: `Rollback package versions`,
	Long:  `Call straight-thaw-versions`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		s, err := newStraight(cmd)
		if err != nil {
			return err
		}
		defer s.Close()

		return s.Rollback(cmd.Context())
	},
}
