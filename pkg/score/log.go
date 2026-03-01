package score

import "github.com/berquerant/grinfo"

type Log struct {
	URL            string  `json:"url" expr:"url"`
	Dir            string  `json:"dir" expr:"dir"`
	LocalCommit    string  `json:"local_commit" expr:"local_commit"`
	LocalTimestamp int64   `json:"local_timestamp" expr:"local_timestamp"`
	LocalDate      string  `json:"local_date" expr:"local_date"`
	LocalRelDate   string  `json:"local_reldate" expr:"local_reldate"`
	DiffDay        int     `json:"diff_day" expr:"diff_day"`
	DiffCommit     int     `json:"diff_commit" expr:"diff_commit"`
	DiffTag        int     `json:"diff_tag" expr:"diff_tag"`
	Score          float64 `json:"score" expr:"score"`
}

func NewLog(x *grinfo.Log) *Log {
	return &Log{
		URL:            x.URL,
		Dir:            x.Dir,
		LocalCommit:    x.Local.Hash,
		LocalTimestamp: x.Local.Author.Timestamp,
		LocalDate:      x.Local.Author.Date,
		LocalRelDate:   x.Local.Author.RelDate,
		DiffDay:        x.Remote.TimeDiffFromLocal.Day,
		DiffCommit:     x.Diff.Commit.Count,
		DiffTag:        x.Diff.Tag.Count,
	}
}
