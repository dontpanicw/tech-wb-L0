package main

import (
	"context"
	"flag"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"log"
	"os"
	"tech-wb-L0/backend/internal/api/handlers"
	"tech-wb-L0/backend/internal/kafka"
	"tech-wb-L0/backend/internal/repository/postgres"
	"tech-wb-L0/backend/internal/service/order"
	pkgHttp "tech-wb-L0/backend/pkg/http"
)

func main() {
	//МИГРАЦИИ
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using system environment variables")
	}

	dsn := os.Getenv("DATABASE_URL")

	m, err := migrate.New("file://migrations", dsn)

	if err != nil {
		log.Fatalf("unable to create migrations instance: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("unable to run migrations: %v", err)
	}
	//if err := m.Down(); err != nil && err != migrate.ErrNoChange {
	//	log.Fatalf("unable to run migrations: %v", err)
	//}
	log.Println("migrations successfully created")

	//Repository

	orderRepo, err := postgres.NewOrderStorage(dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	err = orderRepo.RecoveryCache(context.Background())
	if err != nil {
		log.Fatalf("failed to recovery cache from DB %v", err)
	}
	orderRepo.RangeCache()

	//Kafka
	consumer := kafka.NewConsumer(orderRepo)
	go func() {
		consumer.StartConsumer()
	}()

	//Service
	orderService, err := order.NewOrderService(orderRepo)
	if err != nil {
		log.Fatalf("failed to create order service: %v", err)
	}

	_ = orderService

	handler, err := handlers.NewHandler(orderService)
	if err != nil {
		log.Fatalf("failed to create handler: %v", err)
	}

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8000", "http://127.0.0.1:8000"}, // фронтенд
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
	}))

	handler.WithAuthHandlers(r)

	addr := flag.String("addr", ":8080", "server service address")

	flag.Parse()

	log.Printf("Listening on %s", *addr)
	if err := pkgHttp.CreateAndRunServer(r, *addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
