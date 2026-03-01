package main

import (
	"context"
	"os"

	"github.com/berquerant/emacs-straight-stride/pkg/emacs"
	"github.com/berquerant/emacs-straight-stride/pkg/logx"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "stride",
	Short: `Help manage updates for Emacs packages installed via straight.el`,
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		debug, _ := cmd.Flags().GetBool("debug")
		logx.Setup(os.Stderr, debug)
		return nil
	},
	RunE: func(cmd *cobra.Command, _ []string) error {
		return cmd.Help()
	},
}

func init() {
	rootCmd.PersistentFlags().Bool("debug", false, "enable debug log")
	rootCmd.PersistentFlags().StringP("emacs", "e", "emacs", "emacs binary")
	rootCmd.PersistentFlags().StringP("init", "i", os.Getenv("EMACSD"), "init directory; default: $EMACSD")
}

func newEmacsClient(cmd *cobra.Command) *emacs.Client {
	bin, _ := cmd.Flags().GetString("emacs")
	initDir, _ := cmd.Flags().GetString("init")
	return emacs.NewClient(bin, initDir)
}

func main() {
	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		os.Exit(1)
	}
}
