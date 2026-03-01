package main

import (
	"encoding/json"
	"fmt"
	"iter"
	"log/slog"
	"slices"

	"github.com/berquerant/emacs-straight-stride/pkg/grinfox"
	"github.com/berquerant/emacs-straight-stride/pkg/logx"
	"github.com/berquerant/grinfo"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(infoCmd)
}

func addFetchInfoFlags(cmd *cobra.Command) {
	cmd.Flags().IntP("grinfo-worker", "w", 8, "grinfo worker num")
	cmd.Flags().DurationP("minimum-release-age", "m", 0, "ignore commits/tags newer than this")
	_ = cmd.Flags().MarkHidden("minimum-release-age")
}

func fetchInfo(cmd *cobra.Command) (iter.Seq[*grinfo.Result], error) {
	workerNum, _ := cmd.Flags().GetInt("grinfo-worker")
	minimumReleaseAge, _ := cmd.Flags().GetDuration("minimum-release-age")
	info, err := readInfo(cmd)
	if err != nil {
		return nil, err
	}
	slog.Info("FetchInfo")
	return grinfox.NewFetcher(workerNum, minimumReleaseAge).Fetch(cmd.Context(), slices.Values(info.Directories)), nil
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: `Display git info`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		logs, err := fetchInfo(cmd)
		if err != nil {
			return err
		}
		for x := range logs {
			if err := x.Err; err != nil {
				slog.Warn("failed to fetch grinfo", logx.Err(err))
				continue
			}
			b, err := json.Marshal(x)
			if err != nil {
				slog.Warn("failed to marshal grinfo", logx.Err(err))
				continue
			}
			if _, err := fmt.Printf("%s\n", b); err != nil {
				return err
			}
		}
		return nil
	},
}
