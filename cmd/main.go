package main

import (
	"chat/internal/database"
	"chat/internal/handler"
	"chat/internal/repository"
	"chat/internal/server"
	"chat/internal/service"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		logrus.Fatalf("Error loading .env file. %s", err.Error())
	}

	port := os.Getenv("API_SERVER_PORT")
	currentDB := os.Getenv("CURRENT_DB")

	r, dbType, err := getRepository(currentDB)
	if err != nil {
		logrus.Fatalf("Error occurred while initializing the repository: %s", err.Error())
	}

	s, err := service.NewTodoService(dbType, r)
	if err != nil {
		return
	}
	handler := handler.NewHandler(s)
	routes := handler.InitRoutes(dbType)

	server := &server.WebSocketServer{}
	go func() {
		if err := server.Run(port, routes); err != nil {
			logrus.Fatalf("Error occurred while running WebSocket server: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := server.Shutdown(); err != nil {
		logrus.Fatalf("Error occurred while shutting down WebSocket server: %s", err.Error())
	}
}

func getRepository(currentDB string) (*repository.Repository, string, error) {
	var dbType string
	var db interface{}
	var err error

	switch currentDB {
	case "postgres":
		dbType = repository.PostgresDB
		db, err = initializePostgresDB()
		if err != nil {
			return nil, "", fmt.Errorf("failed to initialize Postgres database: %s", err.Error())
		}

	default:
		return nil, "", fmt.Errorf("unsupported database type: %s", currentDB)
	}
	repo, err := repository.NewRepository(dbType, db)
	if err != nil {
		return nil, "", fmt.Errorf("failed to initialize repository: %s", err.Error())
	}
	return repo, dbType, nil
}

func initializePostgresDB() (*sql.DB, error) {
	if os.Getenv("POSTGRES_HOST") == "" || os.Getenv("POSTGRES_PORT") == "" || os.Getenv("POSTGRES_USER") == "" || os.Getenv("POSTGRES_PASSWORD") == "" || os.Getenv("POSTGRES_DB") == "" || os.Getenv("POSTGRES_SSL_MODE") == "" {
		return nil, fmt.Errorf("some of the required environment variables are not set")
	}

	return database.NewPostgresDB(database.PostgresDB{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DB"),
		SSLMode:  os.Getenv("POSTGRES_SSL_MODE"),
	})
}
