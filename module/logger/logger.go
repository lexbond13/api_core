package logger

import (
	"fmt"
	"github.com/lexbond13/api_core/config"
	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	Log ILogger
)

type ILogger interface {
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(err error)
	Debug(args ...interface{})
	Fatal(args ...interface{})
}

type Logger struct {
	logger *zerolog.Logger
}

type fileHook struct {
	out *os.File
}

func (l *Logger) Info(args ...interface{}) {
	l.logger.Info().Msg(fmt.Sprint(args...))
}

func (l *Logger) Warn(args ...interface{}) {
	l.logger.Warn().Msg(fmt.Sprint(args...))
}

func (l *Logger) Error(err error) {
	l.logger.Error().Msg(err.Error() + ". Stacktrace: " + fmt.Sprintf("%+v\n", err))
	// sent error to sentry
	sentry.CaptureException(err)
}

func (l *Logger) Debug(args ...interface{}) {
	l.logger.Debug().Msg(fmt.Sprint(args...))
}

func (l *Logger) Fatal(args ...interface{}) {
	l.logger.Fatal().Msg(fmt.Sprint(args...))
}

func (f *fileHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if level != zerolog.NoLevel {
		message := fmt.Sprintf("%s|%s|%s \n", time.Now().Format("2006-01-02"), level.String(), msg)
		f.out.Write([]byte(message))
	}
}

func Init(config *config.Logger, isDebug bool) error {
	logLevel := zerolog.InfoLevel
	if isDebug {
		logLevel = zerolog.DebugLevel
	}

	zerolog.SetGlobalLevel(logLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.Logger{}

	// Send logs to file
	if config.FileConfig != nil && config.FileConfig.Path != "" {
		hook, err := newFileHook(config.FileConfig)
		if err != nil {
			return errors.Wrap(err, "fail create logger file hook")
		}
		logger = logger.Hook(hook)
	}

	// Send logs to sentry
	if config.SentryConfig != nil && config.SentryConfig.DSN != "" {
		var err error
		err = NewSentry(config.SentryConfig, isDebug)
		if err != nil {
			return err
		}
	}

	logger.With().Timestamp().Logger()
	Log = &Logger{logger: &logger}

	return nil
}

// create hook for writing logs to file
func newFileHook(config *config.FileConfig) (*fileHook, error) {
	perm, err := strconv.ParseUint(config.Perm, 8, 64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse permissions" )
	}

	err = os.MkdirAll(config.Path, os.FileMode(dirPerm(perm)))
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare directory for log file")
	}

	filePath := config.Path + string(filepath.Separator) + time.Now().Format("2006-01-02") + ".log"

	out, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.FileMode(perm))
	if err != nil {
		return nil, errors.Wrap(err, "failed to open log file")
	}

	return &fileHook{out: out}, nil
}

// dirPerm adds executable bit for triads with read bit: 0640 -> 0750
func dirPerm(perm uint64) uint64 {
	for i := uint64(0); i < 3; i++ {
		// Thanks gofmt, nobody needs spaces in this, sure %)
		if perm&(1<<(i*3+2)) > 0 {
			perm |= 1 << (i * 3)
		}
	}
	return perm
}
