package main

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(commitCmd)
}

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: `Commit package versions`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		s, err := newStraight(cmd)
		if err != nil {
			return err
		}
		defer s.Close()

		return s.Commit(cmd.Context())
	},
}
