package main

import (
	"flag"
	"github.com/go-chi/chi/v5"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
	"test_nanimai/backend/internal/repository/postgres"
	//tables "test_nanimai/backend/pkg/db"
	pkgHttp "test_nanimai/backend/pkg/http"
)

func main() {
	addr := flag.String("addr", ":8080", "handler service address")
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using system environment variables")
	}

	dsn := os.Getenv("DATABASE_URL")

	m, err := migrate.New(
		"file://migrations",
		dsn,
	)
	if err != nil {
		log.Fatalf("failed to initialize migrations: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to initialize migrations: %v", err)
	}
	log.Println("migrations initialized")

	//USERS
	UserRepo, err := postgres.NewUserPostgresStorageUser(dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	_ = UserRepo

	//UserService := service.NewUserService(UserRepo)
	//
	//UserHandlers := handler.NewUser(UserService)

	r := chi.NewRouter()

	//UserHandlers.WithUserHandlers(r)

	log.Printf("Listening on %s", *addr)
	if err := pkgHttp.CreateAndRunServer(r, *addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
