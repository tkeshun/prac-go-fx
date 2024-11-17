// mini-appをあまり低レイヤ関数を使わずに書く
package main

import (
	"app-custom-logger/logger"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

// fxがプロバイダーに要求する条件
// 1. 他のプロバイダーの返り値を引数として取る
// 2. 返り値が依存関係として提供可能であること。他の関数や、fx.Invokeによって依存関係として要求されている場合登録可能

// fx.Invokeは、fx.Provideで登録された依存関係を利用する関数を指定する. 条件は以下
// 1. 引数がすべてfx.Provideで提供される型であること
// 2. 返り値を持たない
// 3. エラー処理は内部で行う。Invokeに指定された関数はエラーをかえせないので、内部で処理する

func main() {
	// シグナル監視用のコンテキストを作成、SIGINT, SIGTERMを監視する。シグナルを検知するとctx.Done()が通知される。
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// ログ設定。共通でJSON形式で標準出力に出力するように設定
	logger.Initialize()

	// fx アプリケーションの起動
	app := fx.New(
		fx.Provide(NewHTTPServer),  // HTTPサーバーの依存関係を提供。この中で依存関係を解決する
		fx.Invoke(StartHTTPServer), // サーバーのライフサイクル（起動/停止）フックを登録
		fx.WithLogger(func() fxevent.Logger {
			return logger.NewJSONLogger()
		}),
	)

	// アプリケーションの実行と終了処理
	// app.Start(ctx)でfxがアプリケーションを開始
	// 依存関係の初期化（fx.Provide）とライフサイクルのOnStartフックを実行
	if err := app.Start(ctx); err != nil {
		logger.Logger.Error("Failed to start application", slog.String("error", err.Error()))
		os.Exit(1)
	}

	<-ctx.Done() // シグナルを受け取るまでブロック、シグナルを受信したら次の処理へ
	slog.Info("Shutting down application...")
	// app.Stop(context.Background())でアプリケーションの停止を開始
	// ライフサイクルのOnStopフックを実行し、HTTPサーバーを安全にシャットダウン
	if err := app.Stop(context.Background()); err != nil {
		logger.Logger.Error("Failed to stop application", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

// HTTPサーバーの生成
func NewHTTPServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		if _, err := io.Copy(w, r.Body); err != nil {
			fmt.Fprintln(os.Stderr, "failed to handle request", err)
		}
	})

	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
}

// サーバーのライフサイクル管理
func StartHTTPServer(lc fx.Lifecycle, srv *http.Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				logger.Logger.Info(fmt.Sprintf("Starting HTTP server at %s", srv.Addr))
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					slog.Error("Server failed", slog.String("error", err.Error()))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Logger.Info("Stopping HTTP server")
			// タイムアウト付きでサーバーをシャットダウン
			shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			return srv.Shutdown(shutdownCtx)
		},
	})
}
