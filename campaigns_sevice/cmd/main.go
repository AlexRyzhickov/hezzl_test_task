package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-redis/redis/v8"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"hezzl_test_task/internal/config"
	"hezzl_test_task/internal/handler"
	repository "hezzl_test_task/internal/repository"
	"hezzl_test_task/internal/utils"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"hezzl_test_task/internal/service"
)

func connectDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DBConn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func connectRedis(cfg *config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: "",
		DB:       0,
	})
	return rdb
}

type Handler interface {
	Method() string
	Path() string
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

func registerHandler(router chi.Router, handler Handler) {
	router.Method(handler.Method(), handler.Path(), handler)
}

func connectionsClosedForServer(server *http.Server, logger *zerolog.Logger) chan struct{} {
	connectionsClosed := make(chan struct{})
	go func() {
		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, os.Interrupt)
		defer signal.Stop(shutdown)
		<-shutdown

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()
		logger.Info().Msg("Closing connections")
		if err := server.Shutdown(ctx); err != nil {
			logger.Error().Err(err).Msg("")
		}
		close(connectionsClosed)
	}()
	return connectionsClosed
}

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	defer nc.Close()

	l := utils.InitLogger(nc)
	logger := zerolog.New(l /*os.Stdout*/).With().Timestamp().Logger()

	cfg, err := config.New()
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	db, err := connectDB(cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect database")
	}

	c := connectRedis(cfg)
	r := &repository.RedisCache{Client: c}
	s := service.NewService(db, r)

	router := chi.NewRouter()
	router.Use(middleware.RequestLogger(&handler.LogFormatter{Logger: &logger}))
	router.Use(middleware.Recoverer)
	router.Use(cors.AllowAll().Handler)

	router.Group(func(router chi.Router) {
		registerHandler(router, handler.NewCreateItemHandler(s, &logger))
		registerHandler(router, handler.NewReadItemsHandler(s, &logger))
		registerHandler(router, handler.NewListItemsHandler(s, &logger))
		registerHandler(router, handler.NewUpdateItemHandler(s, &logger))
		registerHandler(router, handler.NewDeleteItemHandler(s, &logger))
	})

	addr := fmt.Sprintf(":%s", cfg.Port)
	server := http.Server{
		Addr:    addr,
		Handler: router,
	}

	connectionsClosed := connectionsClosedForServer(&server, &logger)
	logger.Info().Msg("Server is listening on " + addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Error().Err(err).Msg("")
	}
	<-connectionsClosed
}
