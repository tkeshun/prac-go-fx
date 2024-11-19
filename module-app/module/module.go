package module

import (
	"module-app/infra"
	"module-app/user"

	"go.uber.org/fx"
)

// Module はユーザーモジュールを定義します。
var Module = fx.Options(
	fx.Provide(infra.NewRepository), // 永続化層の依存を提供
	fx.Provide(user.NewService),     // サービス層の依存を提供
)
