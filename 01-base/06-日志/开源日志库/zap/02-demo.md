

/demo/service/logger_service.go
```golang
package service

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"sync"
)

var (
	httpLogger    *zap.Logger
	consoleLogger *zap.Logger

	output        = map[string]*os.File{}
	loggerMu      sync.Mutex
)

func GetHttpLoggerInstance() *zap.Logger {
	if httpLogger == nil {
		loggerMu.Lock()
		defer loggerMu.Unlock()

		// GetHttpLogConfig 自行实现
		logConfig := GetHttpLogConfig()

		level := zapcore.DebugLevel
		level.Set(logConfig.Level)

	    encoder := getEncoder()
	    writeSyncer := getLogWriter(logConfig)

	    if logConfig.ToStdout {
			writeSyncer = zapcore.NewMultiWriteSyncer(writeSyncer, os.Stdout)
		}

	    core := zapcore.NewCore(encoder, writeSyncer, zap.NewAtomicLevelAt(level))
		//httpLogger = zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel))
	    httpLogger = zap.New(core, zap.AddCaller())
	}

	return httpLogger
}

func GetConsoleLoggerInstance() *zap.Logger {
	if consoleLogger == nil {
		loggerMu.Lock()
		defer loggerMu.Unlock()

		// GetHttpLogConfig 自行实现
		logConfig := GetConsoleLogConfig()

		level := zapcore.DebugLevel
		level.Set(logConfig.Level)

	    encoder := getEncoder()
	    writeSyncer := getLogWriter(logConfig)

	    if logConfig.ToStdout {
			writeSyncer = zapcore.NewMultiWriteSyncer(writeSyncer, os.Stdout)
		}

	    core := zapcore.NewCore(encoder, writeSyncer, zap.NewAtomicLevelAt(level))
		consoleLogger = zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel))
	    // consoleLogger = zap.New(core, zap.AddCaller())
	}

	return consoleLogger
}

func getEncoder() zapcore.Encoder {
    encoderConfig := zap.NewProductionEncoderConfig()
    encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
    return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(logConfig *LogConfig) zapcore.WriteSyncer {
    lumberJackLogger := &lumberjack.Logger{
        Filename:   logConfig.FilePath,
        MaxSize:    logConfig.MaxSize,
        MaxBackups: logConfig.MaxBackups,
        MaxAge:     logConfig.MaxAge,
        Compress:   logConfig.Compress,
    }
    return zapcore.AddSync(lumberJackLogger)
}

type LogConfig struct {
	FilePath          string
	Level             string
	MaxSize           int
	MaxAge            int
	MaxBackups        int
	OutputProbability float32 //日志输出概率(小数)
	ToStdout          bool    //是否将日志输出到标准输出中
	Compress		  bool    //是否压缩
}
```

use
```golang

msg := "aaaabbb"
service.GetHttpLoggerInstance().Info(msg)
service.GetHttpLoggerInstance().Error(msg)

service.GetConsoleLoggerInstance().Error(msg)
```