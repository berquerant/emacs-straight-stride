package main

import (
	"errors"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pruneCmd)
	pruneCmd.AddCommand(pruneCacheCmd, pruneRepoCmd)
	pruneRepoCmd.Flags().Bool("all", false, "prune all repos")
}

var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: `Prune resources`,
}

var pruneCacheCmd = &cobra.Command{
	Use:   "cache",
	Short: `Prune caches`,
	Long:  `Prune straight build caches`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		s, err := newStraight(cmd)
		if err != nil {
			return err
		}
		defer s.Close()

		return s.PruneCache(cmd.Context())
	},
}

var pruneRepoCmd = &cobra.Command{
	Use:   "repo [PATTERN]",
	Short: `Prune repos`,
	Long:  `Prune straight repositories`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := newStraight(cmd)
		if err != nil {
			return err
		}
		defer s.Close()

		all, _ := cmd.Flags().GetBool("all")
		if !all && len(args) == 0 {
			return errors.New("--all or PATTERN is required")
		}
		if all {
			return s.PruneRepo(cmd.Context(), "")
		}
		return s.PruneRepo(cmd.Context(), args[0])
	},
}
