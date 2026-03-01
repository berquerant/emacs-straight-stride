package emacs

import (
	"context"
	"path/filepath"

	"github.com/berquerant/emacs-straight-stride/pkg/execx"
)

type Client struct {
	bin     string
	initDir string
}

func NewClient(bin, initDir string) *Client {
	return &Client{
		bin:     bin,
		initDir: initDir,
	}
}

func (e *Client) newCmd(ctx context.Context, args ...string) *execx.Cmd {
	return execx.NewCmd(ctx, e.bin, append([]string{
		"--batch",
		"--quick",
		"--init-directory", e.initDir,
		"--load", filepath.Join(e.initDir, "init.el"),
	}, args...)...)
}

func (e *Client) Output(ctx context.Context, args ...string) (string, error) {
	return e.newCmd(ctx, args...).Output()
}

func (e *Client) Run(ctx context.Context, args ...string) error {
	return e.newCmd(ctx, args...).Run()
}
