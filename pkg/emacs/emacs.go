package emacs

import (
	"context"
	"os"

	"github.com/berquerant/emacs-straight-stride/pkg/execx"
)

type Client struct {
	bin  string
	args []string
}

func NewClient(bin string, args ...string) *Client {
	return &Client{
		bin:  bin,
		args: args,
	}
}

func (e *Client) newCmd(ctx context.Context, args ...string) *execx.Cmd {
	xs := e.args
	xs = append(xs, "--batch", "--quick")
	xs = append(xs, args...)
	return execx.NewCmd(ctx, e.bin, xs...)
}

func (e *Client) Output(ctx context.Context, args ...string) (string, error) {
	return e.newCmd(ctx, args...).Output()
}

func (e *Client) Run(ctx context.Context, args ...string) error {
	return e.newCmd(ctx, args...).Run()
}

func (e *Client) Exec(ctx context.Context, args ...string) error {
	c := execx.NewCmd(ctx, e.bin, append(e.args, args...)...)
	c.Env = os.Environ()
	return c.Exec()
}
