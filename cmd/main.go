package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todo-hexagonal/internal/adapters/primary/web/handlers"
	"todo-hexagonal/internal/adapters/secondary/cache"
	"todo-hexagonal/internal/adapters/secondary/repository"
	"todo-hexagonal/internal/config"
	"todo-hexagonal/internal/core/domain"

	m "todo-hexagonal/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.MustLoadConfig()

	db := connectToDB(cfg)

	db.AutoMigrate(&domain.User{})

	redisCache, err := cache.NewRedisCache(cfg.RedisHost, cfg.RedisPort, cfg.RedisPassword)
	if err != nil {
		panic(err)
	}

	_ = repository.NewDB(db, redisCache)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	killSig := make(chan os.Signal, 1)

	signal.Notify(killSig, syscall.SIGTERM)

	port := cfg.Port
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: initRoutes(),
	}

	go func() {
		err := srv.ListenAndServe()

		if errors.Is(err, http.ErrServerClosed) {
			fmt.Println("Server closed!")
		} else if err != nil {
			fmt.Printf("Error starting server, %s\n", err)
			os.Exit(1)
		}
	}()

	logger.Info(fmt.Sprintf("Server started successfully at port: %s", port))

	<-killSig

	logger.Info("Server shutting down!")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failed!", slog.Any("err", err))
	}

	logger.Info("Server shutdown successfully!")

}

func initRoutes() *chi.Mux {

	router := chi.NewRouter()

	fileServer := http.FileServer(http.Dir("./static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	router.Group(func(r chi.Router) {
		r.Use(middleware.Logger, m.CSPMiddleware, m.TextHTMLMiddleware)

		r.Get("/", handlers.NewHomeHandler().ServeHTTP)
	})

	return router
}

func connectToDB(cfg *config.Config) *gorm.DB {
	host := cfg.DatabaseHost
	port := cfg.DatabasePort
	user := cfg.DatabaseUser
	password := cfg.DatabasePassword
	dbname := cfg.DatabaseName

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}
