package grinfox

import (
	"context"
	"iter"
	"time"

	"github.com/berquerant/grinfo"
)

type Fetcher struct {
	workerNum         int
	minimumReleaseAge time.Duration
}

func NewFetcher(workerNum int, minimumReleaseAge time.Duration) *Fetcher {
	if workerNum < 1 {
		workerNum = 1
	}
	return &Fetcher{
		workerNum:         workerNum,
		minimumReleaseAge: minimumReleaseAge,
	}
}

func (f *Fetcher) Fetch(ctx context.Context, dirs iter.Seq[string]) iter.Seq[*grinfo.Result] {
	w := grinfo.NewWorker(f.workerNum, 100, f.minimumReleaseAge)
	return w.All(ctx, dirs)
}
