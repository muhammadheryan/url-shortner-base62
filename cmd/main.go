package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/muhammadheryan/url-shortner-base62/application/url"
	"github.com/muhammadheryan/url-shortner-base62/cmd/config"
	_ "github.com/muhammadheryan/url-shortner-base62/docs"
	urlRepo "github.com/muhammadheryan/url-shortner-base62/repository/url"
	"github.com/muhammadheryan/url-shortner-base62/transport"
)

// @title URL Shortener API
// @version 1.0
// @description A URL shortener service with Base62 encoding
// @host localhost:8080
// @BasePath /
func main() {
	// Load configuration from environment variables
	cfg := config.Load()

	log.Printf("Starting server in %s environment", cfg.Environment)

	// Connect to database
	db, err := sqlx.Connect("mysql", cfg.GetDSN())
	if err != nil {
		log.Fatal("err connect db ", err)
	}

	// Set database connection pool settings
	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	// Initialize application layers
	URLRepo := urlRepo.NewURLRepository(db)
	URLApp := url.NewURLApplication(URLRepo)
	httpTransport := transport.NewTransport(URLApp)

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      httpTransport,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	log.Printf("HTTP server running on port %s", cfg.Server.Port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalln("failed server ", err)
	}
}
