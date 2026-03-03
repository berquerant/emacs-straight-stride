package straight

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/berquerant/emacs-straight-stride/pkg/emacs"
	"github.com/berquerant/emacs-straight-stride/pkg/temp"
)

type Info struct {
	Meta        Meta              `json:"meta"`
	Profile     map[string]string `json:"profile"`
	Packages    []string          `json:"packages"`
	Directories []string          `json:"directories"`
}

type Meta struct {
	BaseDir     string `json:"base_dir"`
	RepoDir     string `json:"repo_dir"`
	BuildDir    string `json:"build_dir"`
	BuildCache  string `json:"build_cache"`
	ModifiedDir string `json:"modified_dir"`
}

//go:embed straight.el
var script []byte

type Client struct {
	c          *emacs.Client
	tempDir    *temp.Dir
	scriptPath string
}

func NewClient(c *emacs.Client) (*Client, error) {
	x := &Client{
		c: c,
	}
	if err := x.init(); err != nil {
		return nil, err
	}
	return x, nil
}

func (c *Client) init() error {
	d, err := temp.NewDir("stride/straight")
	if err != nil {
		return err
	}
	c.tempDir = d
	c.scriptPath = d.Join("straight.el")
	return os.WriteFile(c.scriptPath, script, 0644)
}

func (c *Client) Close() error { return c.tempDir.Close() }

func (c *Client) output(ctx context.Context, script string) (string, error) {
	return c.c.Output(ctx, "--load", c.scriptPath, "--eval", script)
}

func (c *Client) run(ctx context.Context, script string) error {
	return c.c.Run(ctx, "--load", c.scriptPath, "--eval", script)
}

func (c *Client) ReadInfo(ctx context.Context) (*Info, error) {
	out, err := c.output(ctx, "(my-straight-write-info)")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to read info", err)
	}
	var x Info
	if err := json.Unmarshal([]byte(out), &x); err != nil {
		return nil, fmt.Errorf("%w: failed to read info", err)
	}
	return &x, nil
}

func (c *Client) Update(ctx context.Context, pkgs ...string) error {
	e := ""
	if len(pkgs) == 0 {
		e = "(my-straight-update nil)"
	} else {
		xs := make([]string, len(pkgs))
		for i, p := range pkgs {
			xs[i] = fmt.Sprintf(`"%s"`, p)
		}
		e = fmt.Sprintf(`(my-straight-update '(%s))`, strings.Join(xs, " "))
	}
	if err := c.run(ctx, e); err != nil {
		return fmt.Errorf("%w: failed to update %v", err, pkgs)
	}
	return nil
}

func (c *Client) Commit(ctx context.Context) error {
	if err := c.run(ctx, "(my-straight-commit)"); err != nil {
		return fmt.Errorf("%w: failed to commit", err)
	}
	return nil
}

func (c *Client) Rollback(ctx context.Context) error {
	if err := c.run(ctx, "(my-straight-rollback)"); err != nil {
		return fmt.Errorf("%w: failed to rollback", err)
	}
	return nil
}

func (c *Client) PruneCache(ctx context.Context) error {
	info, err := c.ReadInfo(ctx)
	if err != nil {
		return fmt.Errorf("%w: failed to prune cache", err)
	}
	d, err := temp.NewDir("stride/straight/prune/cache")
	if err != nil {
		return fmt.Errorf("%w: failed to prune cache", err)
	}
	if err := os.Rename(info.Meta.BuildDir, d.Join("build")); err != nil {
		return fmt.Errorf("%w: failed to prune cache", err)
	}
	if err := os.Rename(info.Meta.BuildCache, d.Join("build-cache.el")); err != nil {
		return fmt.Errorf("%w: failed to prune cache", err)
	}
	if err := os.Rename(info.Meta.ModifiedDir, d.Join("modified")); err != nil {
		return fmt.Errorf("%w: failed to prune cache", err)
	}
	return nil
}

func (c *Client) PruneRepo(ctx context.Context, pattern string) error {
	info, err := c.ReadInfo(ctx)
	if err != nil {
		return fmt.Errorf("%w: failed to prune repo, pattern = %s", err, pattern)
	}
	d, err := temp.NewDir("stride/straight/prune/repo")
	if err != nil {
		return fmt.Errorf("%w: failed to prune repo, pattern = %s", err, pattern)
	}

	if pattern == "" {
		if err := os.Rename(info.Meta.RepoDir, d.Join("repos")); err != nil {
			return fmt.Errorf("%w: failed to prune repo, pattern = %s, dir = %s", err, pattern, info.Meta.RepoDir)
		}
		return nil
	}

	r, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("%w: failed to prune repo, pattern = %s", err, pattern)
	}
	for _, x := range info.Directories {
		b := filepath.Base(x)
		if r.MatchString(b) {
			if err := os.Rename(x, d.Join(b)); err != nil {
				return fmt.Errorf("%w: failed to prune repo, pattern = %s, dir = %s", err, pattern, x)
			}
		}
	}
	return nil
}
