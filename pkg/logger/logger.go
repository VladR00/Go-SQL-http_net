package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func Init() {
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stdout, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: "02.01.2006 15:04:05",
		}),
	))
}
