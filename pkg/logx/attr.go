package logx

import "log/slog"

func Err(err error) any {
	return slog.String("err", err.Error())
}
