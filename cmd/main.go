package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aredoff/reagate/internal/config"
	"github.com/aredoff/reagate/internal/log"
	"github.com/aredoff/reagate/internal/server"
)

func main() {
	configPath := flag.String("config", "/etc/webmon/config.json", "path to json configuration file")
	flag.Parse()

	err := config.LoadOrCreatePersistentConfig(*configPath)
	if err != nil {
		panic(fmt.Sprintf("Couldn't load config: %v", err))
	}

	log.Init()
	srv := server.New()

	go func() {
		srv.Serve(fmt.Sprintf("%s:%d", config.Config.GetString("host"), config.Config.GetInt("port")))
	}()

	log.Info(fmt.Sprintf("Start server on %s:%d", config.Config.GetString("host"), config.Config.GetInt("port")))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down the server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
	log.Info("Server gracefully stopped")
}
