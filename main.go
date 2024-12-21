package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
)

var pool *pgxpool.Pool

// Struct untuk memetakan data dari tabel
type User struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Umur       int    `json:"umur"`
	FotoProfil string `json:"foto_profil"`
	ScoreAwal  int    `json:"score_awal"`
	ScoreAkhir int    `json:"score_akhir"`
	Score      int    `json:"score"`
}

func main() {
	// URL koneksi PostgreSQL dari Railway
	databaseURL := "postgresql://postgres:DeuTfmDFDffpvIwgYQczTLEDZOLnofqV@autorack.proxy.rlwy.net:35693/railway"

	// Membuat koneksi pool
	var err error
	pool, err = pgxpool.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	// Tes koneksi ke database
	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}
	fmt.Println("Successfully connected to Railway PostgreSQL database")

	// Menjalankan HTTP server
	http.HandleFunc("/users", getUsersHandler)
	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Handler untuk mengambil data dari PostgreSQL
func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := pool.Query(context.Background(), "SELECT * FROM users")

 // Ganti nama_tabel dengan nama tabel Anda
	if err != nil {
		http.Error(w, "Failed to execute query", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Umur, &user.FotoProfil, &user.ScoreAwal, &user.ScoreAkhir, &user.Score)
		if err != nil {
			http.Error(w, "Failed to scan row", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	// Mengirim respons sebagai JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
