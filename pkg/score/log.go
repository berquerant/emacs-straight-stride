package score

import (
	"path/filepath"

	"github.com/berquerant/grinfo"
)

type Log struct {
	URL                 string  `json:"url" expr:"url"`
	Dir                 string  `json:"dir" expr:"dir"`
	Basename            string  `json:"basename" expr:"basename"`
	LocalCommit         string  `json:"local_commit" expr:"local_commit"`
	LocalTimestamp      int64   `json:"local_timestamp" expr:"local_timestamp"`
	LocalDate           string  `json:"local_date" expr:"local_date"`
	LocalRelDate        string  `json:"local_reldate" expr:"local_reldate"`
	DiffDay             int     `json:"diff_day" expr:"diff_day"`
	DiffCommit          int     `json:"diff_commit" expr:"diff_commit"`
	DiffTag             int     `json:"diff_tag" expr:"diff_tag"`
	RemoteCommitDiffDay int     `json:"remote_commit_diff_day" expr:"remote_commit_diff_day"`
	RemoteTagDiffDay    int     `json:"remote_tag_diff_day" expr:"remote_tag_diff_day"`
	Score               float64 `json:"score" expr:"score"`
}

func NewLog(x *grinfo.Log) *Log {
	p := &Log{
		URL:                 x.URL,
		Dir:                 x.Dir,
		Basename:            filepath.Base(x.Dir),
		LocalCommit:         x.Local.Hash,
		LocalTimestamp:      x.Local.Author.Timestamp,
		LocalDate:           x.Local.Author.Date,
		LocalRelDate:        x.Local.Author.RelDate,
		DiffDay:             x.Remote.TimeDiffFromLocal.Day,
		DiffCommit:          x.Diff.Commit.Count,
		DiffTag:             x.Diff.Tag.Count,
		RemoteCommitDiffDay: x.Remote.TimeDiffToNow.Day,
	}
	if v := x.RemoteTag; v != nil {
		p.RemoteTagDiffDay = v.TimeDiffToNow.Day
	}
	return p
}
