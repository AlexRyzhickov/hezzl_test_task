package main

import (
	"context"
	"fmt"
	"hezzl_test_task/internal/config"
	"hezzl_test_task/internal/handler"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"hezzl_test_task/internal/service"
	//"hezzl_test_task/internal/subscribe"
)

func connectDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DBConn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

type Handler interface {
	Method() string
	Path() string
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

func registerHandler(router chi.Router, handler Handler) {
	router.Method(handler.Method(), handler.Path(), handler)
}

func connectionsClosedForServer(server *http.Server) chan struct{} {
	connectionsClosed := make(chan struct{})
	go func() {
		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, os.Interrupt)
		defer signal.Stop(shutdown)
		<-shutdown

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()
		log.Println("Closing connections")
		if err := server.Shutdown(ctx); err != nil {
			log.Println(err)
		}
		close(connectionsClosed)
	}()
	return connectionsClosed
}

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	db, err := connectDB(cfg)
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	/*nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("failed to connect nats streaming", err)
	}
	defer nc.Close()

	_, err = nc.Subscribe("foo", func(m *nats.Msg) {
		log.Printf("Received a message: %s\n", string(m.Data))
		if err := subscribe.ProcessOrder(db, m); err != nil {
			log.Println(err)
		}
	})

	if err != nil {
		log.Fatal("failed to subscribe", err)
	}*/

	/*c := cache.InitializeMemoryCache()
	if err = fillCacheFromDatabase(c, db); err != nil {
		log.Fatal(err)
	}*/

	s := service.NewService(db)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(cors.AllowAll().Handler)

	router.Group(func(router chi.Router) {
		registerHandler(router, &handler.CreateItemHandler{Service: s})
		//registerHandler(router, &handler.ReadItemHandler{Service: s})
		registerHandler(router, &handler.ListItemsHandler{Service: s})
		registerHandler(router, &handler.UpdateItemHandler{Service: s})
		registerHandler(router, &handler.DeleteItemHandler{Service: s})
	})

	addr := fmt.Sprintf(":%s", cfg.Port)
	server := http.Server{
		Addr:    addr,
		Handler: router,
	}

	connectionsClosed := connectionsClosedForServer(&server)
	log.Println("Server is listening on " + addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Println(err)
	}
	<-connectionsClosed
}
