package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/4lerman/pm_service/cmd/pm_service/api"
	"github.com/4lerman/pm_service/internal/config"
	"github.com/4lerman/pm_service/pkg/db"
)

func main() {
	db, err := db.NewPSQLStorage(&db.DbConfig{
		Host:     config.Envs.DBAddress,
		User:     config.Envs.DBUser,
		Password: config.Envs.DBPassword,
		Port:     config.Envs.DBPort,
		Dbname:   config.Envs.DBName,
	})

	if err != nil {
		log.Fatal("Db init error", err)
	}

	defer db.Close()

	initStorage(db)

	server := api.NewAPIServer(fmt.Sprint(":", config.Envs.Port), db)

	if err := server.Run(); err != nil {
		log.Fatal("Error when running server: ", err)
	}
}

func initStorage(db *sql.DB) {
	if err := db.Ping(); err != nil {
		log.Fatal("Db connection error", err)
	}

	log.Println("Db connected successfully!")
}
