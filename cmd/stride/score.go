package main

import (
	"encoding/json"
	"fmt"
	"iter"
	"log/slog"
	"math"
	"slices"

	"github.com/berquerant/emacs-straight-stride/pkg/logx"
	"github.com/berquerant/emacs-straight-stride/pkg/score"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(scoreCmd)
	addScoreFlags(scoreCmd)
}

func addScoreFlags(cmd *cobra.Command) {
	addFetchInfoFlags(cmd)
	cmd.Flags().IntP("count", "n", -1, "head count; negative number means infinity")
	cmd.Flags().StringP("calc", "c", score.DefaultCalculatorFormula, "expr of expr-lang; formula to calculate stale score; float value is expected")
	cmd.Flags().StringP("filter", "f", score.DefaultFilterFormula, "expr of expr-lang; formula to filter logs; bool value is expected")
}

func calcScore(cmd *cobra.Command) (iter.Seq[*score.Log], error) {
	calcExpr, _ := cmd.Flags().GetString("calc")
	c, err := score.NewCalculator(calcExpr)
	if err != nil {
		return nil, err
	}

	filterExpr, _ := cmd.Flags().GetString("filter")
	f, err := score.NewFilter(filterExpr)
	if err != nil {
		return nil, err
	}

	logs, err := fetchInfo(cmd)
	if err != nil {
		return nil, err
	}

	return func(yield func(*score.Log) bool) {
		slog.Info("CalcScore")
		for x := range logs {
			if err := x.Err; err != nil {
				slog.Warn("failed to fetch grinfo log", logx.Err(err))
				continue
			}
			logger := slog.With(
				slog.String("url", x.Log.URL),
				slog.String("dir", x.Log.Dir),
			)
			v, err := c.Calculate(x.Log)
			if err != nil {
				logger.Warn("failed to calculate score", logx.Err(err))
				continue
			}
			logger = logger.With(slog.Float64("score", v.Score))
			b, err := f.Select(v)
			if err != nil {
				logger.Warn("failed to filter log", logx.Err(err))
				continue
			}
			if !b {
				logger.Debug("exclude log by filter")
				continue
			}
			if !yield(v) {
				return
			}
		}
	}, nil
}

var scoreCmd = &cobra.Command{
	Use:   "score",
	Short: `Display stale score`,
	Long: `Calculate and display the stale score for each package order by score desc.

The following attributes are available for calculating the stale score (float):
- url (string), repository url
- dir (string), directory of the local repository
- basename (string), basename of the directory of the local repository
- local_commit (string), commit hash of the local repository
- local_timestamp (int), timestamp of the commit
- local_date (string), date of the commit
- local_reldate (string), relative date of the commit
- diff_day (int), time diff between the commit and the latest remote commit
- diff_commit (int), number of commits between the commit and the latest remote commit
- diff_tag (int), number of tags between the commit and the latest remote commit
- remote_commit_diff_day (int), time diff between the latest remote commit and now
- remote_tag_diff_day (int, optional), time diff between the latest tag and now; 0 if not exists

Further attributes available for log filtering:
- score (float)`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		logs, err := calcScore(cmd)
		if err != nil {
			return err
		}
		xs := slices.Collect(logs)
		// sort by score desc
		slices.SortStableFunc(xs, func(a, b *score.Log) int {
			x, y := a.Score, b.Score
			switch {
			case x < y:
				return 1
			case x > y:
				return -1
			default:
				return 0
			}
		})

		maxCount, _ := cmd.Flags().GetInt("count")
		if maxCount < 0 {
			maxCount = math.MaxInt
		}
		var count int
		for _, x := range xs {
			if count >= maxCount {
				return nil
			}
			b, err := json.Marshal(x)
			if err != nil {
				slog.Warn("failed to marshal score", logx.Err(err))
				continue
			}
			if _, err := fmt.Printf("%s\n", b); err != nil {
				return err
			}
			count++
		}
		return nil
	},
}
