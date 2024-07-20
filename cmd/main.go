package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/joho/godotenv"
)

// Response structure
type Response struct {
	Message string `json:"message"`
}

type user struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
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

// getUserHandler fetches a user from the database
func getUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		if email == "" {
			http.Error(w, "missing email parameter", http.StatusBadRequest)
			return
		}

		var u user
		query := "SELECT id, name, email, password_hash FROM users WHERE email = ?"
		err := db.QueryRow(query, email).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "user not found", http.StatusNotFound)
			} else {
				http.Error(w, "server error", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(u)
	}
}

func main() {
	env := os.Getenv("GO_ENV")
	if env != "production" {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Error loading .env file: %s\n", err.Error())
		}
	}

	log.Println("connecting to database...")
	db, err := newMysqlDB()
	if err != nil {
		log.Fatalf("Could not connect to database: %s\n", err.Error())
	}
	defer db.Close()

	log.Println("applying schema migrations...")
	if err := applySchemaMigrations(db); err != nil {
		log.Fatalf("Could not apply schema migrations: %s\n", err.Error())
	}

	log.Println("applying data migrations...")
	if err := applyDataMigrations(db); err != nil {
		log.Fatalf("Could not apply data migrations: %s\n", err.Error())
	}

	http.HandleFunc("/api/hello", helloHandler)
	http.HandleFunc("/api/goodbye", goodbyeHandler)
	http.HandleFunc("/api/user", getUserHandler(db))

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
		"host":     os.Getenv("DB_HOST"),
		"port":     os.Getenv("DB_PORT"),
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config["user"], config["password"], config["host"], config["port"], config["name"])
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
