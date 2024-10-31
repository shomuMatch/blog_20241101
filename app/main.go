package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

type Person struct {
	gorm.Model
	Name string
	Age  int
}

var dbClient *gorm.DB

func operations(tx *gorm.DB) error {

	if err := tx.Create(&Person{Name: "Fujii", Age: 18}).Error; err != nil {
		return fmt.Errorf("first: %w", err)
	}

	var person Person

	if err := tx.First(&person, "name = ?", "Fujii").Error; err != nil {
		return fmt.Errorf("second: %w", err)
	}

	if err := tx.Model(&person).Update("Name", "Fuji").Error; err != nil {
		return fmt.Errorf("third: %w", err)
	}

	if err := tx.Delete(&person).Error; err != nil {
		return fmt.Errorf("fourth: %w", err)
	}
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	tx := dbClient.Begin()
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if err := operations(tx); err != nil {
		tx.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Internal Server Error: %s\n", err.Error())))
		return
	}
	tx.Commit()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK\n"))
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	db, err := gorm.Open("postgres", "user=postgres port=5432 host=extension_test.db password=password sslmode=disable")
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	dbClient = db
	defer dbClient.Close()
	dbClient.AutoMigrate(&Person{})

	http.HandleFunc("/api", handler)

	server := &http.Server{
		Addr:    ":8888",
		Handler: nil,
	}
	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()
	fmt.Println("start receiving at :8888")
	log.Fatal(server.ListenAndServe())
}
