package logging

import "github.com/rs/zerolog"

type AsynqZerologLogger struct {
	Logger *zerolog.Logger
}

func (l *AsynqZerologLogger) Debug(args ...interface{}) {
	l.Logger.Debug().Msgf("%v", args...)
}

func (l *AsynqZerologLogger) Info(args ...interface{}) {
	l.Logger.Info().Msgf("%v", args...)
}

func (l *AsynqZerologLogger) Warn(args ...interface{}) {
	l.Logger.Warn().Msgf("%v", args...)
}

func (l *AsynqZerologLogger) Error(args ...interface{}) {
	l.Logger.Error().Msgf("%v", args...)
}

func (l *AsynqZerologLogger) Fatal(args ...interface{}) {
	l.Logger.Fatal().Msgf("%v", args...)
}
