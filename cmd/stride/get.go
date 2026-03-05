package main

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/berquerant/emacs-straight-stride/pkg/straight"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(getCmd)
}

func newStraight(cmd *cobra.Command) (*straight.Client, error) {
	e, err := newEmacsClient(cmd)
	if err != nil {
		return nil, err
	}
	return straight.NewClient(e)
}

func readInfo(cmd *cobra.Command) (*straight.Info, error) {
	s, err := newStraight(cmd)
	if err != nil {
		return nil, err
	}
	defer s.Close()
	slog.Info("ReadInfo")
	return s.ReadInfo(cmd.Context())
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: `Display resources`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		info, err := readInfo(cmd)
		if err != nil {
			return err
		}
		b, err := json.Marshal(info)
		if err != nil {
			return err
		}
		if _, err := fmt.Printf("%s\n", b); err != nil {
			return err
		}
		return nil
	},
}
