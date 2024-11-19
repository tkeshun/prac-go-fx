package infra

import (
	"context"
	"errors"
	"module-app/user"
)

// DatabaseRepository はユーザー情報を永続化する具体的な実装です。
type DatabaseRepository struct{}

// NewRepository はRepositoryの具体実装を生成します。
func NewRepository() user.Repository {
	return &DatabaseRepository{}
}

// Save はユーザー情報を保存します。
func (r *DatabaseRepository) Save(ctx context.Context, username string) error {
	// 実際のデータベース保存処理（ここではスタブ）
	return nil
}

// Find はユーザー情報を取得します。
func (r *DatabaseRepository) Find(ctx context.Context, id string) (string, error) {
	// 実際のデータベース取得処理（ここではスタブ）
	if id == "1" {
		return "found user", nil
	}
	return "", errors.New("user not found")
}
