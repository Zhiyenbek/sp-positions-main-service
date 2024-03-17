package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/Zhiyenbek/sp-positions-main-service/config"
	handler "github.com/Zhiyenbek/sp-positions-main-service/internal/handler/http"
	"github.com/Zhiyenbek/sp-positions-main-service/internal/repository"
	"github.com/Zhiyenbek/sp-positions-main-service/internal/repository/connection"
	"github.com/Zhiyenbek/sp-positions-main-service/internal/service"
	"go.uber.org/zap"
)

func Run() error {
	logger, _ := zap.NewDevelopment(zap.AddStacktrace(zap.PanicLevel))

	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	cfg, err := config.New()
	if err != nil {
		sugar.Errorf("error while defining config %v", err)
		return err
	}
	db, err := connection.NewPostgresDB(cfg.DB)
	if err != nil {
		sugar.Errorf("error while creating database: %v", err)
		return err
	}
	defer db.Close()
	repos := repository.New(db, cfg, sugar)
	services := service.New(repos, sugar, cfg)
	handlers := handler.New(services, sugar, cfg)

	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Println("Couldn't get port. Using config port instead")
		port = strconv.Itoa(cfg.App.Port)

	}

	srv := http.Server{
		Addr:    ":" + port,
		Handler: handlers.InitRoutes(),
	}
	errChan := make(chan error, 1)
	go func(errChan chan<- error) {
		sugar.Infof("server on port: %d have started", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			sugar.Error(err)
			errChan <- err
		}
	}(errChan)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-quit:
		log.Println("killing signal was received, shutting down the server")
	case err := <-errChan:
		sugar.Errorf("ERROR: HTTP server error received: %v", err)
	}

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.App.TimeOut)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		sugar.Errorf("WARN: Server forced to shutdown: %v", err)
	}
	return nil

}
