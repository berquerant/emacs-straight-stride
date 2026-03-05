package main

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/berquerant/emacs-straight-stride/pkg/emacs"
	"github.com/berquerant/emacs-straight-stride/pkg/logx"
	"github.com/spf13/cobra"
	"mvdan.cc/sh/v3/shell"
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
	rootCmd.PersistentFlags().StringP("emacs", "e", "emacs", "emacs command")
	rootCmd.PersistentFlags().StringP("init-el", "i", "", "init.el; default: $EMACSD/init.el")
	rootCmd.PersistentFlags().StringP("emacsd", "d", "", "override $EMACSD")
}

func defaultEmacsd(cmd *cobra.Command) string {
	if emacsd, _ := cmd.Flags().GetString("emacsd"); emacsd != "" {
		return emacsd
	}
	return os.Getenv("EMACSD")
}

func newEmacsClient(cmd *cobra.Command) (*emacs.Client, error) {
	command, _ := cmd.Flags().GetString("emacs")
	args, err := shell.Fields(command, os.Getenv)
	if err != nil {
		return nil, err
	}
	if len(args) == 0 {
		return nil, errors.New("no emacs command")
	}
	emacsd := defaultEmacsd(cmd)
	args = append(args, "--init-directory", emacsd, "--load", filepath.Join(emacsd, "init.el"))
	return emacs.NewClient(args[0], args[1:]...), nil
}

func main() {
	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		os.Exit(1)
	}
}
