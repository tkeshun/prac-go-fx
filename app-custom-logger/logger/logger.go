package logger

import (
	"fmt"
	"log/slog"
	"os"

	"go.uber.org/fx/fxevent"
)

// Loggerはアプリケーション全体で使用するグローバルロガー
var Logger *slog.Logger

// Initializeはグローバルロガーを初期化します
func Initialize() {
	Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil)) // JSON形式のロガーを設定
}

// JSONLoggerはfx用のカスタムロガー
type JSONLogger struct{}

// NewJSONLoggerはfx.WithLogger用のロガーを作成します
func NewJSONLogger() *JSONLogger {
	return &JSONLogger{}
}

// LogEventはfxのイベントをJSON形式でログ出力します
func (l *JSONLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		Logger.Info("OnStart hook executing",
			slog.String("caller", e.CallerName),
			slog.String("function", e.FunctionName),
		)
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			Logger.Error("OnStart hook failed",
				slog.String("caller", e.CallerName),
				slog.String("function", e.FunctionName),
				slog.String("error", e.Err.Error()),
				slog.String("runtime", e.Runtime.String()),
			)
		} else {
			Logger.Info("OnStart hook executed",
				slog.String("caller", e.CallerName),
				slog.String("function", e.FunctionName),
				slog.String("runtime", e.Runtime.String()),
			)
		}
	case *fxevent.OnStopExecuting:
		Logger.Info("OnStop hook executing",
			slog.String("caller", e.CallerName),
			slog.String("function", e.FunctionName),
		)
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			Logger.Error("OnStop hook failed",
				slog.String("caller", e.CallerName),
				slog.String("function", e.FunctionName),
				slog.String("error", e.Err.Error()),
				slog.String("runtime", e.Runtime.String()),
			)
		} else {
			Logger.Info("OnStop hook executed",
				slog.String("caller", e.CallerName),
				slog.String("function", e.FunctionName),
				slog.String("runtime", e.Runtime.String()),
			)
		}
	case *fxevent.Supplied:
		if e.Err != nil {
			Logger.Error("Supplied failed",
				slog.String("type", e.TypeName),
				slog.String("module", e.ModuleName),
				slog.String("error", e.Err.Error()),
			)
		} else {
			Logger.Info("Supplied",
				slog.String("type", e.TypeName),
				slog.String("module", e.ModuleName),
			)
		}
	case *fxevent.Provided:
		for _, rtype := range e.OutputTypeNames {
			Logger.Info("Provided",
				slog.String("constructor", e.ConstructorName),
				slog.String("module", e.ModuleName),
				slog.String("type", rtype),
			)
		}
		if e.Err != nil {
			Logger.Error("Provide failed",
				slog.String("constructor", e.ConstructorName),
				slog.String("module", e.ModuleName),
				slog.String("error", e.Err.Error()),
			)
		}
	case *fxevent.Invoking:
		Logger.Info("Invoking",
			slog.String("function", e.FunctionName),
			slog.String("module", e.ModuleName),
		)
	case *fxevent.Invoked:
		if e.Err != nil {
			Logger.Error("Invoked failed",
				slog.String("function", e.FunctionName),
				slog.String("module", e.ModuleName),
				slog.String("error", e.Err.Error()),
			)
		}
	case *fxevent.Stopping:
		Logger.Info("Stopping",
			slog.String("signal", e.Signal.String()),
		)
	case *fxevent.Stopped:
		if e.Err != nil {
			Logger.Error("Stopped with error",
				slog.String("error", e.Err.Error()),
			)
		} else {
			Logger.Info("Stopped")
		}
	case *fxevent.RollingBack:
		Logger.Error("Rolling back due to start failure",
			slog.String("error", e.StartErr.Error()),
		)
	case *fxevent.RolledBack:
		if e.Err != nil {
			Logger.Error("Rollback failed",
				slog.String("error", e.Err.Error()),
			)
		} else {
			Logger.Info("Rolled back")
		}
	case *fxevent.Started:
		if e.Err != nil {
			Logger.Error("Start failed",
				slog.String("error", e.Err.Error()),
			)
		} else {
			Logger.Info("Started")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			Logger.Error("Logger initialization failed",
				slog.String("function", e.ConstructorName),
				slog.String("error", e.Err.Error()),
			)
		} else {
			Logger.Info("Logger initialized",
				slog.String("function", e.ConstructorName),
			)
		}
	default:
		Logger.Info("Unhandled event",
			slog.String("event", fmt.Sprintf("%T", event)),
		)
	}
}
