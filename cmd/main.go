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
	"todo-hexagonal/internal/core/services"
	m "todo-hexagonal/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	userService    *services.UserService
	todoService    *services.TodoService
	sessionStore   *sessions.CookieStore
	cfg            *config.Config
	authMiddleware *m.AuthMiddleware
)

func main() {
	cfg = config.MustLoadConfig()

	db := connectToDB(cfg)

	redisCache, err := cache.NewRedisCache(cfg.RedisHost, cfg.RedisPort, cfg.RedisPassword)
	if err != nil {
		panic(err)
	}

	sessionStore = sessions.NewCookieStore([]byte(cfg.SessionCookieName))

	store := repository.NewDB(db, redisCache)

	userService = services.NewUserService(services.NewUserServiceParams{
		Repo: store,
	})
	todoService = services.NewTodoService(services.NewTodoServiceParams{
		Repo: store,
	})

	authMiddleware = m.NewAuthMiddleware(m.NewAuthMiddlewareParams{
		SessionStore:      sessionStore,
		SessionCookieName: cfg.SessionCookieName,
		UserService:       userService,
	})

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
		r.Use(middleware.Logger)

		r.Post("/register", handlers.NewPostRegisterHandler(userService).ServeHTTP)

		r.Post("/login", handlers.NewPostLoginHandler(handlers.NewPostLoginHandlerParams{
			UserService:       userService,
			SessionStore:      sessionStore,
			SessionCookieName: cfg.SessionCookieName,
		}).ServeHTTP)
	})

	router.Group(func(r chi.Router) {
		r.Use(middleware.Logger, authMiddleware.RedirectToHomeIfLoggedIn)

		r.Get("/register", handlers.NewRegisterHandler().ServeHTTP)

		r.Get("/login", handlers.NewLoginHandler().ServeHTTP)
	})

	router.Group(func(r chi.Router) {
		r.Use(middleware.Logger, authMiddleware.ValidateSession)

		r.Get("/", handlers.NewHomeHandler(handlers.NewHomeHandlerParams{
			TodoService: todoService,
		}).ServeHTTP)

		r.Post("/logout", handlers.NewPostLogoutHandler(handlers.NewPostLogoutParams{
			SessionStore:      sessionStore,
			SessionCookieName: cfg.SessionCookieName,
		}).ServeHTTP)

		r.Post("/todos", handlers.NewPostAddTodoHandler(handlers.NewPostAddTodoHandlerParams{
			TodoService: todoService,
		}).ServeHTTP)

		r.Delete("/todos/{id}", handlers.NewPostDeleteTodoHandler(handlers.NewPostDeleteTodoHandlerParams{
			TodoService: todoService,
		}).ServeHTTP)
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

	err = db.AutoMigrate(&domain.User{}, &domain.Todo{})
	if err != nil {
		panic(err)
	}

	return db
}
