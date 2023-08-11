package logging

import "github.com/rs/zerolog"

type AsynqZerologLogger struct {
	Logger *zerolog.Logger
}

func (l *AsynqZerologLogger) Debug(args ...any) {
	l.Logger.Debug().Msgf("%v", args...)
}

func (l *AsynqZerologLogger) Info(args ...any) {
	l.Logger.Info().Msgf("%v", args...)
}

func (l *AsynqZerologLogger) Warn(args ...any) {
	l.Logger.Warn().Msgf("%v", args...)
}

func (l *AsynqZerologLogger) Error(args ...any) {
	l.Logger.Error().Msgf("%v", args...)
}

func (l *AsynqZerologLogger) Fatal(args ...any) {
	l.Logger.Fatal().Msgf("%v", args...)
}
