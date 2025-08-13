package main

// @title Balance Service API
// @version 1.0
// @description API для управления балансом, лимитами и резервами средств
// @BasePath /

import (
	"flag"
	"log"
	"net"
	"os"
	_ "test_nanimai/backend/docs"
	balancegrpc "test_nanimai/backend/internal/api/grpc"
	pb "test_nanimai/backend/internal/api/grpc/pb"
	rest "test_nanimai/backend/internal/api/rest"
	"test_nanimai/backend/internal/repository/postgres"
	"test_nanimai/backend/internal/service/balance"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
)

func main() {
	restAddr := flag.String("rest-addr", ":8080", "REST service address")
	grpcAddr := flag.String("grpc-addr", ":9090", "gRPC service address")
	flag.Parse()

	if env := os.Getenv("REST_ADDR"); env != "" {
		*restAddr = env
	}
	if env := os.Getenv("GRPC_ADDR"); env != "" {
		*grpcAddr = env
	}

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

	// Repositories
	balanceRepo, err := postgres.NewBalanceStorage(dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Services
	balanceService := balance.NewBalanceService(balanceRepo)

	// HTTP server (Gin)
	r := gin.Default()
	// API-key middleware
	r.Use(rest.ApiKeyAuthMiddleware(balanceRepo.GetDb()))
	// REST routes
	rest.RegisterRoutes(r, balanceService)
	// Swagger UI (Gin)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// gRPC server
	lis, err := net.Listen("tcp", *grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", *grpcAddr, err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterBalanceServiceServer(grpcServer, balancegrpc.NewBalanceGRPCServer(balanceService))

	errCh := make(chan error, 2)

	go func() {
		log.Printf("REST listening on %s", *restAddr)
		errCh <- r.Run(*restAddr)
	}()

	go func() {
		log.Printf("gRPC listening on %s", *grpcAddr)
		errCh <- grpcServer.Serve(lis)
	}()

	if err := <-errCh; err != nil {
		log.Fatalf("server stopped with error: %v", err)
	}
}
