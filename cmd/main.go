package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
)

// Response structure
type Response struct {
	Message string `json:"message"`
}

// helloHandler responds with a hello message
func helloHandler(w http.ResponseWriter, r *http.Request) {
	secret := os.Getenv("MY_SECRET")
	port := os.Getenv("PORT")
	anotherSecret := os.Getenv("MY_ANOTHER_SECRET")
	response := Response{Message: fmt.Sprintf("Hello, World! from %s, secret is %s, another secret is %s", port, secret, anotherSecret)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// goodbyeHandler responds with a goodbye message
func goodbyeHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{Message: "Goodbye, World! update"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/api/hello", helloHandler)
	http.HandleFunc("/api/goodbye", goodbyeHandler)

	log.Println("connecting to database...")
	db, err := newMysqlDB()
	if err != nil {
		log.Fatalf("Could not connect to database: %s\n", err.Error())
	}

	log.Println("applying schema migrations...")
	if err := applySchemaMigrations(db); err != nil {
		log.Fatalf("Could not apply schema migrations: %s\n", err.Error())
	}

	log.Println("applying data migrations...")
	if err := applyDataMigrations(db); err != nil {
		log.Fatalf("Could not apply data migrations: %s\n", err.Error())
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "4040"
	}
	log.Println("Starting server on port", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}

func newMysqlDB() (*sql.DB, error) {
	config := map[string]string{
		"name":     os.Getenv("DB_NAME"),
		"user":     os.Getenv("DB_USER"),
		"password": os.Getenv("DB_PASSWORD"),
	}

	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:54321)/%s?parseTime=true", config["user"], config["password"], config["name"])
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func applySchemaMigrations(db *sql.DB) error {
	driver, err := mysql.WithInstance(db, &mysql.Config{
		MigrationsTable: "schema_migrations",
	})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations/schema/",
		"mysql", driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func applyDataMigrations(db *sql.DB) error {
	driver, err := mysql.WithInstance(db, &mysql.Config{
		MigrationsTable: "data_migrations",
	})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations/data/",
		"mysql", driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
