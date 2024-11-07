package main

import (
	"context"
	"log"
	"os"

	"github.com/topten1222/hello_sekai/config"
	"github.com/topten1222/hello_sekai/pkg/database"
	"github.com/topten1222/hello_sekai/server"
)

func main() {
	ctx := context.Background()
	cfg := config.LoadConfig(func() string {
		if len(os.Args) < 2 {
			log.Fatal("Error: .env path is required")
		}
		return os.Args[1]
	}())
	db := database.DbConnect(ctx, &cfg)
	defer db.Disconnect(ctx)

	server.Start(ctx, &cfg, db)
}
