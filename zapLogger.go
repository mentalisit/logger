package Mylogger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

type Logger struct {
	ZapLogger *zap.Logger
}

func LoggerZap(botToken string, chatID int64, webhookDS string) *Logger {
	telegramWriter := NewTelegramWriter(botToken, chatID)
	discordWriter := NewDiscordWriter(webhookDS)

	// Определяем имя файла с логами, включающее "log", дату и время
	logFileName := fmt.Sprintf("log\\log_%s.log", time.Now().Format("2006-01-02_15-04-05"))

	// Определяем WriteSyncer для файла
	fileWriteSyncer := zapcore.AddSync(createLogFile(logFileName))

	cfg := zap.Config{
		Encoding:         "console",
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		OutputPaths:      []string{"stdout", logFileName},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}
	cfgNew := cfg.EncoderConfig
	cfgNew.EncodeLevel = zapcore.CapitalLevelEncoder

	logger, err := cfg.Build(
		zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewTee(core, zapcore.NewCore(
				zapcore.NewConsoleEncoder(cfgNew),
				zapcore.AddSync(telegramWriter),
				cfg.Level,
			), zapcore.NewTee(zapcore.NewCore(
				zapcore.NewConsoleEncoder(cfgNew),
				zapcore.AddSync(discordWriter),
				cfg.Level,
			), zapcore.NewCore(
				zapcore.NewConsoleEncoder(cfgNew),
				fileWriteSyncer,
				cfg.Level,
			)))
		}),
		zap.AddCallerSkip(1),
	)

	if err != nil {
		fmt.Printf("Ошибка при создании логгера: %v\n", err)
		return nil
	}

	defer logger.Sync()
	return &Logger{ZapLogger: logger} // LoggerInterface: logger}
}
func LoggerZapDEV() *Logger {
	cfg := zap.Config{
		Encoding:         "console",
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}
	logger, err := cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil
	}

	defer logger.Sync()

	logger.Info("Develop Running")

	return &Logger{ZapLogger: logger}
}

func (l *Logger) ErrorErr(err error) {
	l.ZapLogger.Error("Произошла ошибка", zap.Error(err))
}
func (l *Logger) Debug(s string, fields ...zap.Field) {
	l.ZapLogger.Debug(s, fields...)
}
func (l *Logger) Info(s string, fields ...zap.Field) {
	l.ZapLogger.Info(s, fields...)
}
func (l *Logger) Warn(s string, fields ...zap.Field) {
	l.ZapLogger.Warn(s, fields...)
}
func (l *Logger) Error(s string, fields ...zap.Field) {
	l.ZapLogger.Error(s, fields...)
}
func (l *Logger) Panic(s string, fields ...zap.Field) {
	l.ZapLogger.Panic(s, fields...)
}
func (l *Logger) Fatal(s string, fields ...zap.Field) {
	l.ZapLogger.Fatal(s, fields...)
}

func (l *Logger) InfoStruct(s string, i interface{}, fields ...zap.Field) {
	l.ZapLogger.Info(fmt.Sprintf("%s: %+v \n", s, i), fields...)
}