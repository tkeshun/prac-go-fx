package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"

	"go.uber.org/fx"
)

func main() {
	// New():コンポーネントなしで構築する
	// Run()停止信号を受信するまでブロックし、終了する前にクリーンアップ操作を実行する
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	fx.New(
		fx.Provide(
			NewHTTPServer,
			NewServeMux,
			NewEchoHandler,
		),
		// Provideで生成された依存関係を利用してアプリケーションの開始処理を定義
		fx.Invoke(func(*http.Server) {}),
	).Run() // SIGINTやSIGTERMのシグナルを受け取ると停止しOnStopを呼び出す
}

// リクエストを受信するサーバーを構築（*http.Serverの実装として使われる）
func NewHTTPServer(lc fx.Lifecycle, mux *http.ServeMux) *http.Server {
	// server設定
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			slog.Info(fmt.Sprintf("Starting HTTP server at %s", srv.Addr))
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}

// リクエストを処理する方法を記述（*http.ServeMuxが入る場所で使われる）
type EchoHandler struct{}

func NewEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

func (*EchoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := io.Copy(w, r.Body); err != nil {
		fmt.Fprintln(os.Stderr, "failed to handle request", err)
	}
}

func NewServeMux(echo *EchoHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/echo", echo)
	return mux
}
