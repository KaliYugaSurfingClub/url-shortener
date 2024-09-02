package sl

import "log/slog"

func ErrorAttr(err error) slog.Attr {
	return slog.String("error", err.Error())
}
