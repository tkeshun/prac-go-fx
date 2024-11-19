package main

import (
	"context"
	"log"
	"module-app/module"
	"module-app/user"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		module.Module,     // モジュール化した依存関係をセット
		fx.Invoke(runApp), // 起動する関数をセット
	).Run()
}

func runApp(service user.Service) {
	ctx := context.Background()

	if err := service.SaveUser(ctx, "example user"); err != nil {
		log.Fatalf("Error saving user: %v", err)
	}

	user, err := service.User(ctx, "1")
	if err != nil {
		log.Fatalf("Error getting user: %v", err)
	}

	log.Printf("Retrieved user: %s", user)
}
