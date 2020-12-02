package logger

import (
	"fmt"

	"go.uber.org/zap"
)

// Logger interface
type Logger interface {
	Debug(msg string, fields ...LogFields)
	Error(msg string, fields ...LogFields)
	Fatal(msg string, fields ...LogFields)
	Info(msg string, fields ...LogFields)
	Panic(msg string, fields ...LogFields)
}

// LogFields to define custom log field key and values
type LogFields struct {
	Key   string
	Value interface{}
}

// Instance of the logger
type Instance struct {
	logger *zap.Logger
}

// Config object to configure the logger
type Config struct {
	MessageKey    string
	CallerKey     string
	NameKey       string
	FunctionKey   string
	LevelKey      string
	StackTraceKey string
	TimestampKey  string
	CallerSkip    int
}

// Default returns a default non-JSON logger
func Default() (*Instance, error) {
	l, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	return &Instance{
		logger: l,
	}, nil
}

// New returns a configurable JSON logger
func New(userConfig Config) (*Instance, error) {
	config := zap.NewProductionConfig()
	applyConfig(&config, userConfig)

	l, err := config.Build()
	if err != nil {
		return nil, err
	}

	if userConfig.CallerSkip != 0 {
		return &Instance{
			logger: l.WithOptions(zap.AddCallerSkip(userConfig.CallerSkip)),
		}, nil
	}

	return &Instance{
		logger: l,
	}, nil
}

// Debug logs
func (li *Instance) Debug(msg string, fields ...LogFields) {
	args := grabFields(fields)
	li.logger.Debug(msg, args...)
}

// Error logs
func (li *Instance) Error(msg string, fields ...LogFields) {
	args := grabFields(fields)
	li.logger.Error(msg, args...)
}

// Fatal logs
func (li *Instance) Fatal(msg string, fields ...LogFields) {
	args := grabFields(fields)
	li.logger.Fatal(msg, args...)
}

// Info logs
func (li *Instance) Info(msg string, fields ...LogFields) {
	args := grabFields(fields)
	li.logger.Info(msg, args...)
}

// Panic logs
func (li *Instance) Panic(msg string, fields ...LogFields) {
	args := grabFields(fields)
	li.logger.Panic(msg, args...)
}

func grabFields(fields []LogFields) []zap.Field {
	args := make([]zap.Field, 0)
	for _, field := range fields {
		key := field.Key
		switch v := field.Value.(type) {
		case int:
			args = append(args, zap.Int(key, v))
		case string:
			args = append(args, zap.String(key, v))
		case bool:
			args = append(args, zap.Bool(key, v))
		case float64:
			args = append(args, zap.Float64(key, v))
		case interface{}:
			args = append(args, zap.Any(key, v))
		case []string:
			for i, val := range v {
				args = append(args, zap.String(key+fmt.Sprint(i), val))
			}
		case []int:
			for i, val := range v {
				args = append(args, zap.Int(key+fmt.Sprint(i), val))
			}
		case []bool:
			for i, val := range v {
				args = append(args, zap.Bool(key+fmt.Sprint(i), val))
			}
		case []float64:
			for i, val := range v {
				args = append(args, zap.Float64(key+fmt.Sprint(i), val))
			}
		case []interface{}:
			for i, val := range v {
				args = append(args, zap.Any(key+fmt.Sprint(i), val))
			}
		}
	}
	return args
}

func applyConfig(config *zap.Config, userConfig Config) {
	if userConfig.CallerKey != "" {
		config.EncoderConfig.CallerKey = userConfig.CallerKey
	}

	if userConfig.FunctionKey != "" {
		config.EncoderConfig.FunctionKey = userConfig.FunctionKey
	}

	if userConfig.LevelKey != "" {
		config.EncoderConfig.LevelKey = userConfig.LevelKey
	}

	if userConfig.MessageKey != "" {
		config.EncoderConfig.MessageKey = userConfig.MessageKey
	}

	if userConfig.NameKey != "" {
		config.EncoderConfig.NameKey = userConfig.NameKey
	}

	if userConfig.StackTraceKey != "" {
		config.EncoderConfig.StacktraceKey = userConfig.StackTraceKey
	}

	if userConfig.TimestampKey != "" {
		config.EncoderConfig.TimeKey = userConfig.TimestampKey
	}
}
