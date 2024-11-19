package user

import "context"

// Repository はユーザー情報の永続化インターフェースです。
type Repository interface {
	Save(ctx context.Context, username string) error
	Find(ctx context.Context, id string) (string, error)
}
