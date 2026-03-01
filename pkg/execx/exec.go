package execx

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

type Cmd struct {
	*exec.Cmd
}

var (
	ErrCmd = errors.New("ErrCmd")
)

func NewCmd(ctx context.Context, name string, args ...string) *Cmd {
	cmd := exec.CommandContext(ctx, name, args...)
	slog.Debug("command", slog.String("dir", cmd.Dir), slog.Any("args", cmd.Args))
	return &Cmd{cmd}
}

func (c *Cmd) Output() (string, error) {
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)
	c.Stdout = &stdout
	c.Stderr = &stderr
	if err := c.Cmd.Run(); err != nil {
		return "", fmt.Errorf("%w: command=%v, %s", errors.Join(ErrCmd, err), c.Args, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func (c *Cmd) Run() error {
	c.Cmd.Stdout = os.Stdout
	c.Cmd.Stderr = os.Stderr
	if err := c.Cmd.Run(); err != nil {
		return fmt.Errorf("%w: command=%v", errors.Join(ErrCmd, err), c.Args)
	}
	return nil
}
