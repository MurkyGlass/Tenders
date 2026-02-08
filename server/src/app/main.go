package main

import (
	"main/internal/repositories"
	"main/internal/router"
	handler "main/internal/router/handlers"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	var logger = logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Out = os.Stdout
		logger.Info("Failed to log to file, using default stderr")
	} else {
		logger.Out = file
		defer file.Close()
	}

	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()


	repo := repositories.NewRepository(db)

	handlers := handler.NewHandlers(logger, repo)

	router := router.NewRouter(handlers,db)
	router.Use(handler.LoggingMiddleware(logger))
	router.Use(handler.CorsMiddleware())

	logger.Info("Starting server on :8080")
	logger.Fatal(http.ListenAndServe(":8080", router))
}
